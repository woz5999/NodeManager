package consumer

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg/constants"
	"github.com/woz5999/NodeManager/pkg/node"
	"github.com/woz5999/NodeManager/pkg/queue"
	"github.com/woz5999/NodeManager/pkg/types"
)

// Consumer a consumer worker thread
type Consumer struct {
	Base  *types.Base
	EC2   *ec2.EC2
	Queue *queue.Queue
}

// Start start the worker thread
func (c Consumer) Start(ctx context.Context) error {
	ticker := time.Tick(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker:
				msg, err := c.Queue.Read()
				if err != nil {
					log.Error(err.Error())
					continue
				}

				if msg == nil {
					log.Info("No message in queue")
					continue
				}

				event, err := msg.Body()
				if err != nil {
					log.Error(err.Error())
					continue
				}

				// determine if we care about this event
				if event.LifecycleTransition != constants.InstanceTerminating {
					log.Infof("Received lifecycle transition %s. Ignoring...", string(event.LifecycleTransition))
					err = msg.Delete()
					if err != nil {
						log.Error(err.Error())
					}
					continue
				}

				// create node struct from the ec2 id in the parsed message
				n := &node.Node{
					EC2InstanceID: event.EC2InstanceID,
					EC2:           c.EC2,
				}

				err = n.Drain()
				if err != nil {
					log.Error(err.Error())
					msg.Visibility()
					continue
				}

				err = msg.Delete()
				if err != nil {
					log.Error(err.Error())
				}
			case <-ctx.Done():
				return
			}
		}
	}()
	return nil
}
