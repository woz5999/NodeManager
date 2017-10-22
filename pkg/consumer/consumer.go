package consumer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg/constants"
	"github.com/woz5999/NodeManager/pkg/event"
	"github.com/woz5999/NodeManager/pkg/node"
	"github.com/woz5999/NodeManager/pkg/types"
)

// Consumer a consumer worker thread
type Consumer struct {
	Base *types.Base
	Svc  *sqs.SQS
}

// Start start the worker thread
func (c Consumer) Start(ctx context.Context) error {
	ticker := time.Tick(1 * time.Second)
	go func() {
		for {
			select {
			case <-ticker:
				result, err := c.Svc.ReceiveMessage(&sqs.ReceiveMessageInput{
					AttributeNames: []*string{
						aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
					},
					MessageAttributeNames: []*string{
						aws.String(sqs.QueueAttributeNameAll),
					},
					QueueUrl:            &c.Base.Config.AwsSqsQueueURL,
					MaxNumberOfMessages: aws.Int64(1),
					VisibilityTimeout:   aws.Int64(c.Base.Config.ErrorVisibilityTimeoutSec),
					WaitTimeSeconds:     aws.Int64(0),
				})

				if err != nil {
					log.Error(err.Error())
					break
				}

				if len(result.Messages) == 0 {
					log.Info("Empty Queue. Pausing 5 Seconds")
					time.Sleep(5 * time.Second)
					break
				}

				// process message
				msg := result.Messages[0]
				event := event.Event{}
				err = json.Unmarshal([]byte(*msg.Body), &event)
				if err != nil {
					log.Error(err.Error())
					c.updateVisibility(msg)
					break
				}

				// determine if we care about this event
				if event.LifecycleTransition != constants.InstanceTerminating {
					log.Info("Received lifecycle transition " + string(event.LifecycleTransition) + ". Ignoring...")
					err = c.deleteMessage(msg)
					if err != nil {
						log.Error(err.Error())
					}
					break
				}

				// create node struct from the ec2 id in the parsed message
				n := &node.Node{
					Base:          c.Base,
					EC2InstanceID: event.EC2InstanceID,
				}

				err = n.Drain()
				if err != nil {
					log.Error(err.Error())
					c.updateVisibility(msg)
					break
				}

				err = c.deleteMessage(msg)
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

func (c Consumer) deleteMessage(msg *sqs.Message) error {
	_, err := c.Svc.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &c.Base.Config.AwsSqsQueueURL,
		ReceiptHandle: msg.ReceiptHandle,
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (c Consumer) updateVisibility(msg *sqs.Message) error {
	_, err := c.Svc.ChangeMessageVisibility(&sqs.ChangeMessageVisibilityInput{
		QueueUrl:          &c.Base.Config.AwsSqsQueueURL,
		ReceiptHandle:     msg.ReceiptHandle,
		VisibilityTimeout: aws.Int64(c.Base.Config.ErrorVisibilityTimeoutSec),
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
