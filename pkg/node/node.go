package node

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg/kubectl"
)

// Node an EC2 node that needs to be drained
type Node struct {
	EC2           *ec2.EC2
	EC2InstanceID string
	instance      ec2.Instance
}

// Drain all pods from the node using its aws private hostname
func (n Node) Drain() error {
	hostname, err := n.PrivateHostname()
	if err != nil {
		log.Error(err.Error())
		return err
	}

	log.Infof("Draining node %s", hostname)
	k := &kubectl.Kubectl{}
	err = k.Exec([]string{"drain", hostname, "--force"})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (n Node) initInstance() error {
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
	res, _ := n.EC2.DescribeInstances(params)

	switch {
	case len(res.Reservations[0].Instances) == 1:
		log.Info("Found instance with ID %s", n.EC2InstanceID)
		n.instance = *res.Reservations[0].Instances[0]
	case len(res.Reservations[0].Instances) > 1:
		log.Errorf("Too many instances found with ID %s", n.EC2InstanceID)
		return fmt.Errorf("Too many instances found with ID %s", n.EC2InstanceID)
	default:
		log.Errorf("No instances found with ID %s", n.EC2InstanceID)
		return fmt.Errorf("No instances found with ID %s", n.EC2InstanceID)
	}
	return nil
}

// PrivateHostname the EC2 instance's private hostname
func (n Node) PrivateHostname() (string, error) {
	if &n.instance == nil {
		err := n.initInstance()
		if err != nil {
			log.Error(err.Error())
			return "", err
		}
	}

	if &n.instance == nil {
		log.Errorf("Unknown error retrieving hostname. Instance %s not defined", n.EC2InstanceID)
		return "", fmt.Errorf("Unknown error retrieving hostname. Instance %s not defined", n.EC2InstanceID)
	}
	return *n.instance.PrivateDnsName, nil
}
