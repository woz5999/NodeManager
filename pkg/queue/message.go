package queue

import (
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg/event"
	"github.com/woz5999/NodeManager/pkg/types"
)

// Message an SQS message wrapper
type Message struct {
	Msg  *sqs.Message
	SQS  sqs.SQS
	Base *types.Base
}

// Body return the message event body
func (m Message) Body() (*event.Event, error) {
	// process message
	event := &event.Event{}
	err := json.Unmarshal([]byte(*m.Msg.Body), &event)
	if err != nil {
		log.Error(err.Error())
		m.Delete()
		return nil, err
	}
	return event, nil
}

// Delete the message
func (m Message) Delete() error {
	log.Infof("Deleting mesage %s", m.Msg.MessageId)
	_, err := m.SQS.DeleteMessage(&sqs.DeleteMessageInput{
		QueueUrl:      &m.Base.Config.AwsSqsQueueURL,
		ReceiptHandle: m.Msg.ReceiptHandle,
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

// Visibility set configured error visibility timeout
func (m Message) Visibility() error {
	log.Infof("Updating visibility for mesage %s", m.Msg.MessageId)
	_, err := m.SQS.ChangeMessageVisibility(&sqs.ChangeMessageVisibilityInput{
		QueueUrl:          &m.Base.Config.AwsSqsQueueURL,
		ReceiptHandle:     m.Msg.ReceiptHandle,
		VisibilityTimeout: aws.Int64(m.Base.Config.ErrorVisibilityTimeoutSec),
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}
