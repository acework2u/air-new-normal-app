package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
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

//func EnvMongoURI() string {
//	err := godotenv.Load()
//
//	if err != nil {
//		fmt.Println("Error loading .env file")
//		log.Fatal(err.Error())
//		//return os.Getenv("DB_URL")
//	}
//
//	return os.Getenv("DB_URL")
//}

func HandleRequest(ctx context.Context, iot *AirIot) (newId string, err2 error) {

	if err != nil {
		log.Panicln(err.Error())
	}

	query := bson.M{
		"device_sn": iot.DeviceSn,
		"timestamp": iot.Timestamp,
		"message":   iot.Message,
	}
	defer client.Disconnect(ctx)
	collection := client.Database("air_newnormal").Collection("airs_things")
	cursor, err := collection.InsertOne(ctx, query)

	if err != nil {
		log.Println(err.Error())
		return "", err
	}
	resp := fmt.Sprintf("%s", cursor.InsertedID)

	return resp, nil

}

func main() {
	//airInfo := &AirIot{
	//	DeviceSn:  "2306F01054323",
	//	Timestamp: 1699938683409,
	//	Message:   "Wpx7g3bvwr.skjfdfdokoofdof",
	//}

	//HandleRequest(ctx, airInfo)

	lambda.Start(HandleRequest)

}
