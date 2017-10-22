package types

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/woz5999/NodeManager/pkg/config"
	"k8s.io/client-go/kubernetes"
)

// Base is a struct for embedding
type Base struct {
	AwsSess   *session.Session
	K8sClient *kubernetes.Clientset
	Config    *config.Config
}
