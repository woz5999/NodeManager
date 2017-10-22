package nodeman

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

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
	initialCtx, cancelCtx := context.WithCancel(context.Background())
	log.Info("Initializing global context")

	handleInterrupt(cancelCtx)

	// init sqs queue
	svc := sqs.New(nm.AwsSess)

	// start consumer threads
	for i := 0; i <= nm.Config.ConsumerThreads; i++ {
		consumer := consumer.Consumer{
			Base: nm.Base,
			Svc:  svc,
		}
		err := consumer.Start(initialCtx)
		if err != nil {
			log.Error(err.Error())
			cancelCtx()
		}
	}

	<-initialCtx.Done()
}

func handleInterrupt(cancelCtx context.CancelFunc) {
	// cleanly handle sig int
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		for sig := range c {
			// got sig int
			log.Printf("Caught sig %v..  exiting...", sig)
			log.Info("Cancelling context...")
			cancelCtx()
			log.Info("Waiting a sec to facilitate cleanup..")
			time.Sleep(time.Second)
			log.Info("Exiting")
			os.Exit(1)
		}
	}()
}
