package config

import (
	"context"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var (
	ctx context.Context
)

func ConnectDB(dbUrl string) *mongo.Client {
	log.Println("Connecting to DB")
	log.Println(dbUrl)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbUrl))
	if err != nil {
		log.Fatal(err)
	}

	return client

}

func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {

	if godotenv.Load() != nil {
		dbName := viper.GetString("DB_Name")
		collection := client.Database(dbName).Collection(collectionName)
		return collection
	} else {

		dbName := os.Getenv("DB_Name")
		collection := client.Database(dbName).Collection(collectionName)
		return collection
	}

}
