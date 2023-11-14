package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"
)

type AirIot struct {
	DeviceSn  string    `bson:"deviceSn" json:"deviceSn"`
	Timestamp time.Time `bson:"timestamp" json:"timestamp"`
	Message   string    `bson:"message" json:"message"`
}

var (
	ctx         context.Context
	DBUrl       = EnvMongoURI()
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(DBUrl))
)

func EnvMongoURI() string {
	err := godotenv.Load()

	if err != nil {
		fmt.Println("Error loading .env file")
		log.Fatal(err.Error())
		//return os.Getenv("DB_URL")
	}

	return os.Getenv("DB_URL")
}

func HandleRequest(ctx context.Context, iot *AirIot) {

	if err != nil {
		log.Panicln(err.Error())
	}

	query := bson.M{
		"device_sn": iot.DeviceSn,
		"timestamp": iot.Timestamp,
		"message":   iot.Message,
	}

	collection := client.Database("air_newnormal").Collection("airs_things")
	cursor, err := collection.InsertOne(ctx, query)

	if err != nil {
		log.Println(err.Error())
	}

	log.Println(cursor.InsertedID)
}

func main() {
	airInfo := &AirIot{
		DeviceSn:  "2306F01054323",
		Timestamp: time.Now(),
		Message:   "Wpx7g3bvwr.skjfdfdokoofdof",
	}

	HandleRequest(ctx, airInfo)

}
