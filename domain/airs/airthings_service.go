package airs

import (
	"Airnewnormal/repository"
	"Airnewnormal/utils"
	"context"
	"fmt"
	"log"
	"time"
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
func (s *airThingsService) AirThingsById(sn string, filter *Filter) ([]*AirReport, error) {

	result := []*AirNewNormal{}

	if len(sn) > 12 {

		dbRs, err := s.airRepo.ReadAirIndoorValId(sn, filter.StartAt, filter.EndAt)

		if err != nil {
			return nil, err
		}

		for _, item := range dbRs {
			rs := &AirNewNormal{
				DeviceSn:  item.DeviceSn,
				Message:   readMsg(item.Message),
				Timestamp: item.Timestamp.Local(),
			}

			result = append(result, rs)
		}

	}

	airReport := []*AirReport{}
	if len(result) > 0 {

		for _, k := range result {
			air := &AirReport{
				DeviceSn: k.DeviceSn,
				IndVal:   IndoorVal{Power: k.Message.Power, Temp: k.Message.Temp, RoomTemp: k.Message.RoomTemp},
				TimeAt:   fmt.Sprintf("%v", k.Timestamp),
			}
			airReport = append(airReport, air)
		}

	}

	return airReport, nil
}

func powerVal(p string) int {

	if p == "on" {
		return 1
	}

	return 0
}

func readTime(tm time.Time) time.Time {

	nt := fmt.Sprintf("%v", tm.Local())
	fmt.Println("This tm =", tm)
	fmt.Println("this nt = ", nt)
	t, _ := time.Parse(time.RFC3339, fmt.Sprintf("%v", tm))
	fmt.Println("This t = ", t)
	return t
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
