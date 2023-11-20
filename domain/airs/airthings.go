package airs

import "time"

type AirNewNormal struct {
	DeviceSn  string      `json:"device_sn"`
	Message   *IndoorInfo `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
}

type Filter struct {
	DeviceSn string `json:"device_sn"`
	StartAt  string `json:"start_at"`
	EndAt    string `json:"end_at"`
}

type AirIAQReport struct {
	DeviceSn string  `json:"deviceSn"`
	Temp     float64 `json:"temp"`
	SetTemp  float64 `json:"setTemp"`
	Time     string  `json:"time"`
}

type IndoorInfo struct {
	Power    string `json:"power,omitempty"`
	Mode     string `json:"mode,omitempty"`
	Temp     string `json:"temp,omitempty"`
	RoomTemp string `json:"roomTemp,omitempty"`
	RhSet    string `json:"rhSet,omitempty"`
	RhRoom   string `json:"RhRoom,omitempty"`
	FanSpeed string `json:"fanSpeed,omitempty"`
	Louver   string `json:"louver,omitempty"`
}

type AirService interface {
	AirThings() ([]*AirNewNormal, error)
	AirThingsById(deviceSn string, fil *Filter) ([]*AirNewNormal, error)
}
