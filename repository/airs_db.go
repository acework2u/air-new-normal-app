package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type AirRepositoryDB struct {
	airsThingsCollection *mongo.Collection
	ctx                  context.Context
}

func NewAirRepositoryDB(ctx context.Context, airsThingsCollection *mongo.Collection) AirsRepository {
	return &AirRepositoryDB{
		ctx:                  ctx,
		airsThingsCollection: airsThingsCollection,
	}
}

func (r *AirRepositoryDB) ReadAirIndoorVal() ([]*AirRawData, error) {
	//query := bson.D{{Key: "$match", Value: bson.M{"device_sn": "2306F01054324"}}}

	query := bson.M{}
	//cursor, err := r.airsThingsCollection.Aggregate(r.ctx, query)
	cursor, err := r.airsThingsCollection.Find(r.ctx, query)

	if err != nil {
		return nil, err
	}

	defer cursor.Close(r.ctx)

	devices := []*AirRawData{}

	for cursor.Next(r.ctx) {
		airVal := &AirRawData{}
		ok := cursor.Decode(airVal)
		if ok != nil {
			return nil, ok
		}
		devices = append(devices, airVal)
	}

	return devices, nil
}
func (r *AirRepositoryDB) ReadAirIndoorValId(filter *Filter) ([]*AirRawData, error) {

	stDate := primitive.NewDateTimeFromTime(filter.StartDateAt)
	endDate := primitive.NewDateTimeFromTime(filter.EndDateAt.AddDate(0, 0, 1))

	_ = stDate
	_ = endDate

	//matchStage := bson.D{{"$match", bson.D{{"device_sn", filter.DeviceSN}, {"$and", []bson.M{bson.M{"timestamp": bson.M{"$gte": stDate}}, bson.M{"timestamp": bson.M{"$lt": endDate}}}}}}}
	// match stage
	matchStage := bson.D{{"$match", bson.M{"device_sn": filter.DeviceSN}}}

	if filter.StartDateAt.IsZero() && filter.EndDateAt.IsZero() {
		matchStage = bson.D{{"$match", bson.M{"device_sn": filter.DeviceSN}}}
	}
	if !filter.StartDateAt.IsZero() && filter.EndDateAt.IsZero() {
		matchStage = bson.D{{"$match", bson.M{"device_sn": filter.DeviceSN, "timestamp": bson.M{"$gte": stDate}}}}
	}
	if filter.StartDateAt.IsZero() && !filter.EndDateAt.IsZero() {
		matchStage = bson.D{{"$match", bson.M{"device_sn": filter.DeviceSN, "timestamp": bson.M{"$lt": endDate}}}}
	}
	if !filter.StartDateAt.IsZero() && !filter.EndDateAt.IsZero() {
		matchStage = bson.D{{"$match", bson.M{"device_sn": filter.DeviceSN, "timestamp": bson.M{"$gte": stDate, "$lt": endDate}}}}
	}

	sortStage := bson.D{{"$sort", bson.M{"timestamp": -1}}}

	log.Println("matchStage")
	log.Println(matchStage)

	if filter.Limit == 0 {
		filter.Limit = 1000
	}
	// limit stage

	var cursor *mongo.Cursor
	var ok error

	if filter.Limit > 0 {
		limitStage := bson.D{{"$limit", filter.Limit}}
		cursor, ok = r.airsThingsCollection.Aggregate(r.ctx, mongo.Pipeline{matchStage, sortStage, limitStage})
		if ok != nil {
			return nil, ok
		}
	} else {
		cursor, ok = r.airsThingsCollection.Aggregate(r.ctx, mongo.Pipeline{matchStage, sortStage})
		if ok != nil {
			return nil, ok
		}
	}
	defer cursor.Close(r.ctx)

	// print result
	devices := []*AirRawData{}
	for cursor.Next(r.ctx) {
		airVal := &AirRawData{}
		ok := cursor.Decode(airVal)
		if ok != nil {
			return nil, ok
		}
		devices = append(devices, airVal)
	}

	return devices, nil
}
