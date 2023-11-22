package airs

import "time"

type AirNewNormal struct {
	DeviceSn  string      `json:"device_sn"`
	Message   *IndoorInfo `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
}

type AirReport struct {
	DeviceSn string    `json:"device_sn"`
	IndVal   IndoorVal `json:"ind_val"`
	TimeAt   time.Time `json:"time_at"`
}

type AirInGrafana struct {
	DeviceSn string    `json:"deviceSn"`
	IndVal   IndValue  `json:"indVal"`
	TimeAt   time.Time `json:"timeAt"`
}

type IndValue struct {
	Power    int     `json:"power"`
	Temp     float64 `json:"temp"`
	RoomTemp float64 `json:"roomTemp"`
}

type IndoorVal struct {
	Power    string `json:"power"`
	Temp     string `json:"temp"`
	RoomTemp string `json:"roomTemp"`
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
	AirThingsById(deviceSn string, fil *Filter) ([]*AirReport, error)
	AirThingsById2(deviceSn string, fil *Filter) ([]*AirInGrafana, error)
}
