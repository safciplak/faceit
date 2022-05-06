package awssns

import (
	"context"
	"time"

	"bitbucket.org/faceit/app"

	"bitbucket.org/faceit/internal/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

func New(sess *session.Session, topic string) *SNS {
	return &SNS{snsClient: sns.New(sess), topic: topic}
}

type SNS struct {
	snsClient *sns.SNS
	topic     string
}

func (s *SNS) Publish(ctx context.Context, e events.Event) error {

	if app.IsDEV() {
		return nil
	}

	attrs := make(map[string]*sns.MessageAttributeValue)

	for k, v := range e.Attrs() {
		attrs[k] = &sns.MessageAttributeValue{
			DataType:    aws.String("String"),
			StringValue: aws.String(v),
		}
	}

	_, err := s.snsClient.Publish(&sns.PublishInput{
		TopicArn:          aws.String(s.topic),
		Message:           aws.String(time.Now().String()),
		MessageAttributes: attrs,
	})
	if err != nil {
		return err
	}

	return nil
}
