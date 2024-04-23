package main

import (
	"Airnewnormal/pkg/aws"
	"Airnewnormal/pkg/logger"
	"context"
	"github.com/aws/aws-lambda-go/lambda"
	"go.uber.org/zap"
)

var (
	ctx context.Context
)

func init() {

	log, _ := zap.NewProduction()
	conn, err := aws.New()
	if err != nil {
		log.Fatal("could not set up aws client", zap.Error(err))

	}
	// inject those clients into a context , and start lambda with that context
	ctx = context.Background()

	// inject logger and aws client into context
	ctx = logger.Inject(ctx, log)
	ctx = aws.Inject(ctx, conn)

}

func handler(ctx context.Context, event interface{}) {
	log := logger.GetLoggerFromContext(ctx)
	log.Info("received event lambda", zap.Any("event", event))

	awsServ := aws.GetConnectionFromContext(ctx)

	err := awsServ.SendSQSMessage(ctx, "this message from sqs")
	if err != nil {
		log.Error("could not send sqs ", zap.Error(err))
	}
	err = awsServ.PublishSNSMessage(ctx, "Hello sns")
	if err != nil {
		log.Error("could not send sns ", zap.Error(err))
	}

}
func main() {
	lambda.Start(handler)
}
