package airs

import (
	"Airnewnormal/repository"
	"Airnewnormal/utils"
	"context"
	"fmt"
	"log"
	"strconv"
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

	startAt, _ := time.Parse(time.RFC3339, filter.StartAt)
	finishAt, _ := time.Parse(time.RFC3339, filter.EndAt)

	d := time.Now()
	a := utils.ConvDateDB(d)

	fmt.Println("FInish Add fat = ", d)
	fmt.Println("FInish Add fat a = ", a)

	if len(sn) > 12 {

		queryFilter := repository.Filter{
			DeviceSN:    sn,
			StartDateAt: startAt,
			EndDateAt:   finishAt,
		}

		dbRs, err := s.airRepo.ReadAirIndoorValId(&queryFilter)

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
				TimeAt:   k.Timestamp,
			}
			airReport = append(airReport, air)
		}

	}

	return airReport, nil
}
func (s *airThingsService) AirThingsById2(sn string, filter *Filter) ([]*AirInGrafana, error) {

	result := []*AirNewNormal{}

	startAt, _ := time.Parse(time.RFC3339, filter.StartAt)
	finishAt, _ := time.Parse(time.RFC3339, filter.EndAt)

	log.Println("startAt = ", startAt)
	log.Println("End At -=", finishAt)

	if len(sn) > 12 {

		queryFilter := repository.Filter{
			DeviceSN:    sn,
			StartDateAt: startAt,
			EndDateAt:   finishAt,
		}

		dbRs, err := s.airRepo.ReadAirIndoorValId(&queryFilter)

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

	airReport := []*AirInGrafana{}
	if len(result) > 0 {

		for _, k := range result {

			setTemp, _ := strconv.ParseFloat(k.Message.Temp, 64)
			roomTemp, _ := strconv.ParseFloat(k.Message.RoomTemp, 64)
			//atTime := k.Timestamp.Local().Unix()
			atTime := k.Timestamp.Local()

			air := &AirInGrafana{
				DeviceSn: k.DeviceSn,
				IndVal:   IndValue{Power: powerVal(k.Message.Power), Temp: setTemp, RoomTemp: roomTemp},
				TimeAt:   atTime,
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

func readTime(tm time.Time) string {

	//nt := fmt.Sprintf("%v", tm.Local())
	//fmt.Println("This tm =", tm)
	//fmt.Println("this nt = ", nt)
	//t, _ := time.Parse(time.RFC3339, fmt.Sprintf("%v", tm))
	//fmt.Println("This t = ", t)
	//dateTex := tm.Format("2017.09.07 17:06:06")

	dateTex := fmt.Sprintf("%s", tm)

	return dateTex
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
