package aws

import (
	"context"
)

type AwsKeyType string

var awsKey AwsKeyType = "AWS"

func Inject(ctx context.Context, log *Connection) context.Context {
	return context.WithValue(ctx, awsKey, log)
}

func GetConnectionFromContext(ctx context.Context) *Connection {
	c, _ := ctx.Value(awsKey).(*Connection)
	return c
}
