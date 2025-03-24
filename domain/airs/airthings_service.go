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
func (s *airThingsService) AirThingsById(sn string, filter *Filter) ([]*AirNewNormal, error) {

	startAt, _ := time.Parse(time.RFC3339, filter.StartAt)
	finishAt, _ := time.Parse(time.RFC3339, filter.EndAt)

	queryFilter := repository.Filter{
		DeviceSN:    sn,
		StartDateAt: startAt,
		EndDateAt:   finishAt,
		Limit:       filter.Limit,
	}
	log.Println("queryFilter = ", queryFilter)
	dbRs, err := s.airRepo.ReadAirIndoorValId(&queryFilter)

	if err != nil {
		return nil, err
	}

	airReport := []*AirNewNormal{}
	for _, item := range dbRs {
		air := &AirNewNormal{
			DeviceSn:  item.DeviceSn,
			Message:   readMsg(item.Message),
			Timestamp: item.Timestamp,
		}
		airReport = append(airReport, air)

	}
	return airReport, nil

	//airReport := []*AirReport{}
	/*
		if len(result) > 0 {

			for _, k := range result {

				air := &AirReport{
					DeviceSn: k.DeviceSn,
					IndVal:   IndoorVal{Power: k.Message.Power, Temp: k.Message.Temp, RoomTemp: k.Message.RoomTemp, RhSet: k.Message.RhSet, RhRoom: k.Message.RhRoom},
					TimeAt:   k.Timestamp,
				}
				airReport = append(airReport, air)
			}

		}

		return airReport, nil

	*/
}
func (s *airThingsService) AirThingsById2(sn string, filter *Filter) ([]*AirInGrafana, error) {

	result := []*AirNewNormal{}

	startAt, _ := time.Parse(time.RFC3339, filter.StartAt)
	finishAt, _ := time.Parse(time.RFC3339, filter.EndAt)
	//
	//log.Println("startAt = ", startAt)
	//log.Println("End At -=", finishAt)

	if len(sn) > 12 {

		queryFilter := repository.Filter{
			DeviceSN:    sn,
			StartDateAt: startAt,
			EndDateAt:   finishAt,
			Limit:       filter.Limit,
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
func (s *airThingsService) DeviceThingsById(filter *Filter) ([]*AirInv2000Val, error) {

	deviceSn := filter.DeviceSn
	startAt, _ := time.Parse(time.RFC3339, filter.StartAt)
	endAt, _ := time.Parse(time.RFC3339, filter.EndAt)

	queryFilter := repository.Filter{
		DeviceSN:    deviceSn,
		StartDateAt: startAt,
		EndDateAt:   endAt,
		Limit:       filter.Limit,
	}

	indVal, err := s.airRepo.ReadAirIndoorValId(&queryFilter)

	if err != nil {
		return nil, err
	}

	ac2000 := []*AirInv2000Val{}

	for _, item := range indVal {
		rs := &AirInv2000Val{
			DeviceSn:  item.DeviceSn,
			Message:   readInv2000(item.Message),
			Timestamp: item.Timestamp,
		}

		ac2000 = append(ac2000, rs)
	}

	return ac2000, nil
}
func (s *airThingsService) DeviceThingsOduVal(filter *Filter) ([]*AirOdu3000, error) {
	deviceSn := filter.DeviceSn
	startAt, _ := time.Parse(time.RFC3339, filter.StartAt)
	endAt, _ := time.Parse(time.RFC3339, filter.EndAt)

	queryFilter := repository.Filter{
		DeviceSN:    deviceSn,
		StartDateAt: startAt,
		EndDateAt:   endAt,
		Limit:       filter.Limit,
	}
	ac3000 := []*AirOdu3000{}
	if len(deviceSn) > 12 {

		oduVal, err := s.airRepo.ReadAirIndoorValId(&queryFilter)

		if err != nil {
			return nil, err
		}

		for _, item := range oduVal {
			odu := &AirOdu3000{
				DeviceSN:  item.DeviceSn,
				OduVal:    readOdu3000(item.Message),
				Timestamp: item.Timestamp,
			}

			ac3000 = append(ac3000, odu)
		}

	} // end if

	return ac3000, nil
}

func powerVal(p string) int {

	if p == "on" {
		return 1
	}

	return 0
}

func readTime(tm time.Time) string {

	dateTex := fmt.Sprintf("%s", tm)

	return dateTex
}

func readMsg(msg string) *IndoorInfo {

	airDecode, err := utils.GetClaimsFromToken(msg)
	if err != nil {
		log.Println(err)
		return nil
	}

	//reg1000 := airDecode["data"].(map[string]interface{})["reg1000"].(string)
	//reg2000 := airDecode["data"].(map[string]interface{})["reg2000"].(string)
	//log.Println("reg1000 = ", reg1000)
	//log.Println("reg2000 = ", reg2000)
	//log.Println("reg2000[14:18] = ", reg2000[14:18])
	acValReq := utils.DecodeValAcShadow(airDecode)
	acVal := utils.NewGetAcVal(acValReq)
	ac1000 := acVal.Ac1000()
	acData := (*IndoorInfo)(ac1000)

	return acData
}

func readInv2000(msg string) *IndVal2000 {
	ind2000 := &IndVal2000{}
	airDecode, err := utils.GetClaimsFromToken(msg)
	if err != nil {
		return nil
	}

	reg2000 := airDecode["data"].(map[string]interface{})["reg2000"].(string)
	acVal := utils.NewAcVal("2", reg2000)
	ac2000 := acVal.Ac2000()
	ind2000 = (*IndVal2000)(ac2000)

	return ind2000
}
func readOdu3000(msg string) *OduVal3000 {
	odu3000 := &OduVal3000{}
	airDecode, err := utils.GetClaimsFromToken(msg)
	if err != nil {
		return nil
	}

	reg3000 := airDecode["data"].(map[string]interface{})["reg3000"].(string)
	acVal := utils.NewAcVal("3", reg3000)
	ac3000 := acVal.Ac3000()
	odu3000 = (*OduVal3000)(ac3000)

	return odu3000
}
