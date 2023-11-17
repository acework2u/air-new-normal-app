package airs

import (
	"Airnewnormal/repository"
	"context"
)

type airThingsService struct {
	Ctx     context.Context
	airRepo repository.AirsRepository
}

func NewAirThingsService(airRepo repository.AirsRepository) AirService {

	return &airThingsService{Ctx: context.TODO(), airRepo: airRepo}
}

func (s *airThingsService) AirThings() ([]*AirNewNormal, error) {

	airList, err := s.airRepo.ReadAirIndoorVal()

	if err != nil {
		return nil, err
	}

	airVal := []*AirNewNormal{}

	for _, item := range airList {
		air := &AirNewNormal{
			DeviceSn:  item.DeviceSn,
			Message:   item.Message,
			Timestamp: item.Timestamp,
		}
		airVal = append(airVal, air)
	}

	return airVal, nil
}
