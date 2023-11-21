package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
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
func (r *AirRepositoryDB) ReadAirIndoorValId(deviceSn string, st string, end string) ([]*AirRawData, error) {

	// $match:{device_sn:"2306F01054324",$and:[{timestamp:{$gte:ISODate("2023-11-20")}},{timestamp:{$lte:ISODate("2023-11-21")}}]}
	//	matchStage := bson.D{{"$match", bson.D{{"device_sn", deviceSn}}}}
	//stDate := fmt.Sprintf("ISODate('%s')", st)

	date, _ := time.Parse("2006-01-02", st)
	endAt, _ := time.Parse("2006-01-02", end)

	stDate := primitive.NewDateTimeFromTime(date)
	endDate := primitive.NewDateTimeFromTime(endAt)
	//andStage := bson.D{{"$and", bson.D{{"timestamp", bson.D{{"$gte", stDate}}}}}}
	//matchStage := bson.D{{"$match", bson.D{{"device_sn", deviceSn}}}}
	matchStage := bson.D{{"$match", bson.D{{"device_sn", deviceSn}, {"$and", []bson.M{bson.M{"timestamp": bson.M{"$gte": stDate}}, bson.M{"timestamp": bson.M{"$lte": endDate}}}}}}}
	sortStage := bson.D{{"$sort", bson.M{"timestamp": 1}}}
	//matchStage := bson.D{{"$match", bson.D{{"device_sn", deviceSn}, {"$and", bson.D{{"timestamp", bson.M{"$gte": stDate}}}}}}}
	//matchStage := bson.D{{"$match", bson.D{{"device_sn", ""}, {"$and", bson.D{{"timestamp", bson.D{{"$gte", fmt.Sprintf("ISODate('%v')", st)}}}}}}}}

	//log.Println(matchStage)

	cursor, err := r.airsThingsCollection.Aggregate(r.ctx, mongo.Pipeline{matchStage, sortStage})
	if err != nil {
		return nil, err
	}

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
