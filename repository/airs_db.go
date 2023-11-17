package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
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
func (r *AirRepositoryDB) ReadAirIndoorValId(deviceSn string) ([]*AirRawData, error) {
	return nil, nil
}
