package main

import (
	"context"
	"fmt"
	_ "github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type AirIot struct {
	DeviceSn  string    `bson:"device_sn" json:"device_sn"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Message   string    `bson:"message" json:"message"`
}

var (
	ctx         context.Context
	DBUrl       = os.Getenv("DB_URL")
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(DBUrl))
)

func HandleRequest(ctx context.Context, iot *AirIot) (newId string, err2 error) {

	if err != nil {
		log.Panicln(err.Error())
	}

	query := bson.M{
		"device_sn": iot.DeviceSn,
		"timestamp": iot.Timestamp,
		"message":   iot.Message,
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Panicln("mongo disconnect err: #{err}")
		}
	}()

	collection := client.Database("airs").Collection("airs_shadows")
	cursor, err := collection.InsertOne(ctx, query)

	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	resp := fmt.Sprintf("%s", cursor.InsertedID)

	return resp, nil

}

func main() {
	//lambda.Start(HandleRequest)

}
