package config

import (
	"context"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	ctx context.Context
)

func ConnectDB(dbUrl string) *mongo.Client {

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}

	return client

}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	dbName := viper.GetString("DB_Name")
	collection := client.Database(dbName).Collection(collectionName)
	return collection

}
