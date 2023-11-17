package airs

import "time"

type AirNewNormal struct {
	DeviceSn  string    `json:"device_sn"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

type AirService interface {
	AirThings() ([]*AirNewNormal, error)
}
