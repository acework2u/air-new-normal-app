package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type AirRawData struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	DeviceSn  string             `bson:"device_sn" json:"device_sn"`
	Message   string             `bson:"message" json:"message"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
}

type Filter struct {
	DeviceSN    string
	StartDateAt time.Time
	EndDateAt   time.Time
}

type AirsRepository interface {
	ReadAirIndoorVal() ([]*AirRawData, error)
	ReadAirIndoorValId(deviceSn string) ([]*AirRawData, error)
}
