package main

import (
	"context"
	"fmt"
	"os"
)

func main() {
	ctx := context.Background()
	// Database init
	mc, err := Init(ctx, "DB_URL")
	if err != nil {
		fmt.Printf("mongo error: %v", err)
		os.Exit(1)
	}
	defer mc.Disconnect(ctx)

	//Create a mongo database with the db name
	mongoDB := mc.Database("notification_service")
	//Creat a notification token collection
	tokenCollection := mongoDB.Collection("notificationTokens")

	_ = tokenCollection

}
