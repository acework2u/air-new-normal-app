package airs

import (
	"Airnewnormal/repository"
	"Airnewnormal/utils"
	"context"
	"log"
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
			Message:   readMsg(item.Message),
			Timestamp: item.Timestamp,
		}
		airVal = append(airVal, air)
	}

	return airVal, nil
}
func readMsg(msg string) *IndoorInfo {

	airDecode, err := utils.GetClaimsFromToken(msg)
	if err != nil {
		log.Println(err)
		return nil
	}

	reg1000 := airDecode["data"].(map[string]interface{})["reg1000"].(string)
	acVal := utils.NewGetAcVal(reg1000)
	ac1000 := acVal.Ac1000()
	acData := (*IndoorInfo)(ac1000)

	return acData
}
