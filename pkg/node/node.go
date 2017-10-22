package node

import (
	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg/types"
	"k8s.io/api/core/v1"
	// "k8s.io/client-go/pkg/api/v1"
)

type Node struct {
	EC2InstanceID string
	*types.Base
	Name string
}

func (n Node) runningPods() ([]*v1.Pod, error) {
	log.Info("Getting running pods on Node " + n.Name)

	// list pods
	return nil, nil
}

func (n Node) drain() error {
	log.Info("Draining code " + n.Name)
	// needs a timeout
	return nil
}

func (n Node) leaveCluster() error {
	log.Info("Removing node " + n.Name + " from cluster")
	return nil
}

func (n Node) Terminate() error {
	err = n.drain()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = n.leaveCluster()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
