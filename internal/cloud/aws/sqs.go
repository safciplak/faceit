package aws

import (
	"context"
	"fmt"

	"bitbucket.org/faceit/internal/cloud"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type SQS struct {
	client *sqs.SQS
}

func NewSQS(client *sqs.SQS) SQS {
	return SQS{
		client: client,
	}
}

func (s SQS) Send(ctx context.Context, req *cloud.SendRequest) (string, error) {
	attrs := make(map[string]*sqs.MessageAttributeValue, len(req.Attributes))
	for _, attr := range req.Attributes {
		attrs[attr.Key] = &sqs.MessageAttributeValue{
			StringValue: aws.String(attr.Value),
			DataType:    aws.String(attr.Type),
		}
	}

	res, err := s.client.SendMessageWithContext(ctx, &sqs.SendMessageInput{
		MessageAttributes: attrs,
		MessageBody:       aws.String(req.Body),
		QueueUrl:          aws.String(req.QueueURL),
	})
	if err != nil {
		return "", fmt.Errorf("send: %w", err)
	}

	return *res.MessageId, nil
}
