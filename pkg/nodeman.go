package nodeman

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg/config"
	"github.com/woz5999/NodeManager/pkg/constants"
	"github.com/woz5999/NodeManager/pkg/consumer"
	"github.com/woz5999/NodeManager/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type NodeMan struct {
	*types.Base
}

// NewBase instantiates and returns the base structure used throughout nodeman
func NewBase(config *config.Config) (*types.Base, error) {
	// AWS Session
	awsSess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(config.AwsRegion),
	}))

	// Kubernetes API client.
	k8sClient, err := newK8sClient()
	if err != nil {
		return nil, err
	}

	base := &types.Base{
		AwsSess:   awsSess,
		K8sClient: k8sClient,
		Config:    config,
	}

	return base, nil
}

// NewNodeMan instantiates and returns nodeman.
func NewNodeMan(base *types.Base) (*NodeMan, error) {
	nodeman := &NodeMan{
		Base: base,
	}

	return nodeman, nil
}

func newK8sClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	config.UserAgent = constants.UserAgentBase + constants.Version

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

// Watch watches the SQS Queue for messages to remove nodes.
func (nm *NodeMan) Watch() {
	initialCtx, _ := context.WithCancel(context.Background())
	log.Info("Initializing global context")

	// init sqs queue
	svc := sqs.New(nm.AwsSess)

	// start consumer threads
	for i := 0; i <= nm.Config.ConsumerThreads; i++ {
		consumer := consumer.Consumer{
			Base:           nm.Base,
			Svc:            svc,
			AwsSqsQueueURL: nm.Config.AwsSqsQueueURL,
		}
		err := consumer.Start(initialCtx)
		if err != nil {
			log.Error(err.Error())
		}
	}

	<-initialCtx.Done()
}
