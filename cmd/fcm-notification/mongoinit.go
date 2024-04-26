package main

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Function to connect DB
func Init(ctx context.Context, URL string) (*mongo.Client, error) {
	URL = ""
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(URL))
	if err != nil {
		return nil, err
	}
	// Ping the database to check connection
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}
	fmt.Println("Successfully Connected to the Database")

	return client, nil

}
