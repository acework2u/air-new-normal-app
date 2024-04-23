package aws

import (
	"Airnewnormal/pkg/logger"
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"go.uber.org/zap"
	"os"
)

type Connection struct {
	sqs *sqs.Client
	sns *sns.Client

	queueURL string
	topARN   string
}

func New() (*Connection, error) {

	region := "ap-southeast-1"
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(region))
	if err != nil {
		return nil, err
	}

	return &Connection{
		sqs:      sqs.NewFromConfig(cfg),
		sns:      sns.NewFromConfig(cfg),
		queueURL: os.Getenv("QUEUE_URL"),
		topARN:   os.Getenv("TOPIC_ANR"),
	}, nil
}

func (c *Connection) SendSQSMessage(ctx context.Context, message string) error {
	log := logger.GetLoggerFromContext(ctx)
	output, err := c.sqs.SendMessage(ctx, &sqs.SendMessageInput{
		MessageBody: &message,
		QueueUrl:    &c.queueURL,
	})

	log.Info("send sqs message", zap.Any("output", output))

	return err
}
func (c *Connection) PublishSNSMessage(ctx context.Context, message string) error {
	log := logger.GetLoggerFromContext(ctx)
	output, err := c.sns.Publish(ctx, &sns.PublishInput{
		Message:  &message,
		TopicArn: &c.topARN,
	})

	log.Info("publish sns message", zap.Any("output", output))

	return err
}
