package airs

import "time"

type AirNewNormal struct {
	DeviceSn  string      `json:"device_sn"`
	Message   *IndoorInfo `json:"message"`
	Timestamp time.Time   `json:"timestamp"`
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
}
