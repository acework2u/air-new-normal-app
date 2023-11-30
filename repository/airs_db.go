package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

	matchStage := bson.D{{"$match", bson.D{{"device_sn", filter.DeviceSN}, {"$and", []bson.M{bson.M{"timestamp": bson.M{"$gte": stDate}}, bson.M{"timestamp": bson.M{"$lt": endDate}}}}}}}
	sortStage := bson.D{{"$sort", bson.M{"timestamp": -1}}}

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

	//log.Println(matchStage)

	//cursor, err := r.airsThingsCollection.Aggregate(r.ctx, mongo.Pipeline{matchStage, sortStage})
	//if err != nil {
	//	return nil, err
	//}

	defer cursor.Close(r.ctx)
	//result := bson.M{}
	//cursor.All(r.ctx, &result)
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
