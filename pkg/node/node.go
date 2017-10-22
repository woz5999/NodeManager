package node

import (
	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg/kubectl"
	"github.com/woz5999/NodeManager/pkg/types"
	// "k8s.io/client-go/pkg/api/v1"
)

type Node struct {
	EC2InstanceID string
	*types.Base
	Name    string
	kubectl *kubectl.Kubectl
}

func (n Node) Drain() error {
	log.Info("Draining code " + n.Name)
	err := n.kubectl.Exec([]string{"drain", n.Name, "--force"})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
