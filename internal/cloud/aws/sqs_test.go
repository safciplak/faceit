package aws

import (
	"context"
	"log"
	"testing"

	"bitbucket.org/faceit/internal/cloud"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sqs"
)

type Config struct {
	Address string
	Region  string
	Profile string
	ID      string
	Secret  string
}

func TestQueue(t *testing.T) {
	sess := session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials("ID", "SECRET", ""),
		Region:      aws.String("eu-west-1"),
		// LogLevel:    aws.LogLevel(aws.LogDebug),
	})

	sqsClient := sqs.New(sess)

	sqsService := NewSQS(sqsClient)

	req := cloud.SendRequest{
		QueueURL: "https://sqs.eu-west-1.amazonaws.com/1231321213/queue_url",
		Attributes: []cloud.Attribute{
			{Key: "event", Value: "event", Type: "String"},
			{Key: "params", Value: "params", Type: "String"},
		},
		Body: "test message body",
	}

	id, err := sqsService.Send(context.Background(), &req)

	log.Println(id, err)
}
