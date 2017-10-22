package node

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg/kubectl"
	"github.com/woz5999/NodeManager/pkg/types"
)

// Node an EC2 node that needs to be drained
type Node struct {
	EC2InstanceID string
	*types.Base
	kubectl   *kubectl.Kubectl
	instance  ec2.Instance
	privateIP string
	publicIP  string
}

// Drain all pods from the node using public and/or private ip as node name
func (n Node) Drain() error {
	err := n.initInstance()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	n.drain(*n.instance.PublicIpAddress)
	n.drain(*n.instance.PrivateIpAddress)
	return nil
}

func (n Node) initInstance() error {
	svc := ec2.New(n.Base.AwsSess)

	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id "),
				Values: []*string{
					aws.String(n.EC2InstanceID),
				},
			},
		},
		MaxResults: aws.Int64(1),
	}

	res, _ := svc.DescribeInstances(params)

	for _, i := range res.Reservations[0].Instances {
		n.instance = *i
	}
	return nil
}

func (n Node) drain(node string) error {
	log.Info("Draining node " + node + " with EC2 ID " + n.EC2InstanceID)
	err := n.kubectl.Exec([]string{"drain", node, "--force"})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
