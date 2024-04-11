package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
	_ "github.com/aws/aws-lambda-go/lambda"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

type AirIoTDB struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DeviceSn  string             `bson:"device_sn" json:"device_sn"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Message   string             `bson:"message" json:"message"`
}

var (
	ctx    context.Context
	DB_URL = os.Getenv("DB_URL")
	client *mongo.Client
)

//func envMongoURI() string {
//	err := godotenv.Load()
//	if err != nil {
//		fmt.Println("DB Error loading .env file")
//		log.Panicln(err.Error())
//	}
//	return os.Getenv("DB_URL")
//}

func dbConnect() error {
	var err error
	client, err = mongo.Connect(context.Background(), options.Client().ApplyURI(DB_URL))
	if err != nil {
		// log.Panicln(err.Error())
		return err
	}
	fmt.Println("DB Connected success")

	return nil
}

func init() {
	// dbConnect()
	DB_URL = os.Getenv("DB_URL")
}

// func HandleRequest(ctx context.Context, iot *AirIot) (newId string, err2 error) {
//
//	//if err != nil {
//	//	log.Panicln(err.Error())
//	//}
//
//	query := bson.M{
//		"device_sn": iot.DeviceSn,
//		"timestamp": iot.Timestamp,
//		"message":   iot.Message,
//	}
//	defer func() {
//		if err := client.Disconnect(ctx); err != nil {
//			log.Panicln("mongo disconnect err: #{err}")
//		}
//	}()
//
//	collection := client.Database("airs").Collection("airs_shadows")
//	cursor, err := collection.InsertOne(ctx, query)
//
//	if err != nil {
//		log.Panicln("database Disconnected")
//		log.Println(err.Error())
//		return "", err
//	}
//	resp := fmt.Sprintf("%s", cursor.InsertedID)
//
//	return resp, nil
//
// }

func HandleRequest2(ctx context.Context, iot *AirIot) (newId string, err2 error) {

	if len(DB_URL) == 0 {
		os.Getenv("DB_URL")
	}
	err := dbConnect()

	if err != nil {
		log.Panic(err)
	}

	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Panicln("mongo disconnect err: #{err}")
		}
	}()

	filter := bson.M{"device_sn": string(iot.DeviceSn)}
	query := bson.M{
		"$set": bson.M{"device_sn": iot.DeviceSn, "message": iot.Message, "timestamp": time.Now()},
	}

	upsert := true
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
		Upsert:         &upsert,
	}

	collection := client.Database("airs").Collection("airs_shadows")

	cursor := collection.FindOneAndUpdate(ctx, filter, query, &opt)

	resp := AirIoTDB{}

	if ok := cursor.Decode(&resp); ok != nil {
		return "", ok
	}

	return resp.DeviceSn, nil

}
func main() {
	lambda.Start(HandleRequest2)

	//ac := &AirIot{
	//	DeviceSn:  "2306F01054323",
	//	Timestamp: time.Now(),
	//	Message:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.ewoJInNlcmlhbE51bWJlciI6CSIyMzA2RjAxMDU0MzIzIiwKCSJ3aWZpIjoJewoJCSJzc2lkIjoJIk1vbHRob2xfU1BTIEAyLjRHIiwKCQkiYWlyTmFtZSI6CSJJbmRvb3IwIiwKCQkiYWlyUGFzc3dvcmQiOgkiMDAwMCIsCgkJIm1hY0FkZHJlc3MiOgkiNjA1NUY5N0VENTMwIiwKCQkiaXBBZGRyZXNzIjoJIjE5Mi4xNjguMS4xMDUiLAoJCSJ2ZXJzaW9uIjoJIjguMC4wIgoJfSwKCSJtZXRhIjoJewoJCSJlcnJvcnMiOglbMjYsIDI5LCAyNiwgMjksIDI2LCAyOSwgMCwgMCwgMCwgMF0sCgkJInNlcnZlckNvbm5lY3RlZCI6CXRydWUsCgkJImluQ29ubmVjdGVkIjoJdHJ1ZSwKCQkib3V0Q29ubmVjdGVkIjoJdHJ1ZQoJfSwKCSJkYXRhIjoJewoJCSJyZWcxMDAwIjoJIjAwMDAwMDAwMDAzMjAwODYwMDMyMDA0MTAwMDAwMDA1MDAyMDAwMDAiLAoJCSJyZWcyMDAwIjoJIjAwNDUwMDQ3MDAwMDAwMDAwMDAwMDAxQTAwNDIwQkI4MDAwMDAwRkEiLAoJCSJyZWczMDAwIjoJIjAwNDYwMEZGMDA0NzAwRkYwMDRBMDA0NjAwNTQwMEZGMDAwMTAwNTAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDAwMDEyMDAwMDAwNzgwMDAwMDA3ODAwMDAwMDAwMDAwRTAwMDAwMDAwMDAwMDAwMDAwMEZBMDAxQTAwMDUwMDAwMDAwMDAwMUUwMDI4MDAxOTAwMDAwMDQyMDA0MjAwNDIiLAoJCSJyZWc0MDAwIjoJIjAwODgwMDAwMDAwMDAwMDMwMDE3MDA4MDAwMDAwMDAwMDAwMDAwMDAiCgl9Cn0.IoNSt68HbwrO9fl0eDqEnJHOUj8fIoijvTJ5BVLxogo",
	//}
	//
	//res, err := HandleRequest2(ctx, ac)
	//if err != nil {
	//	fmt.Println(err.Error())
	//}
	//fmt.Println(res)

}
