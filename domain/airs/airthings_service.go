package airs

import "context"

type airThingsService struct {
	Ctx context.Context
}

func NewAirThingsService() AirService {

	return &airThingsService{Ctx: context.TODO()}
}

func (s *airThingsService) AirThings() ([]*AirNewNormal, error) {

	return nil, nil
}
