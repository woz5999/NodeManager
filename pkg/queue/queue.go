package queue

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
	log "github.com/sirupsen/logrus"
	"github.com/woz5999/NodeManager/pkg/types"
)

// Queue an SQS queue wrapper
type Queue struct {
	SQS  sqs.SQS
	Base *types.Base
}

func (q Queue) Read() (*Message, error) {
	result, err := q.SQS.ReceiveMessage(&sqs.ReceiveMessageInput{
		AttributeNames: []*string{
			aws.String(sqs.MessageSystemAttributeNameSentTimestamp),
		},
		MessageAttributeNames: []*string{
			aws.String(sqs.QueueAttributeNameAll),
		},
		QueueUrl:            &q.Base.Config.AwsSqsQueueURL,
		MaxNumberOfMessages: aws.Int64(1),
		VisibilityTimeout:   aws.Int64(q.Base.Config.DefaultVisibilityTimeoutSec),
		WaitTimeSeconds:     aws.Int64(q.Base.Config.QueueWaitTimeSec),
	})

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	if len(result.Messages) < 1 {
		return nil, nil
	}

	m := &Message{
		Msg:  result.Messages[0],
		Base: q.Base,
		SQS:  q.SQS,
	}
	return m, nil
}
