package handler

import (
	"Airnewnormal/domain/airs"
	"Airnewnormal/utils"
	"github.com/gin-gonic/gin"
	"log"
)

type AirsHandler struct {
	airIot airs.AirService
	resp   utils.Response
}

func NewAirsHandler(airIot airs.AirService) AirsHandler {

	return AirsHandler{airIot: airIot, resp: utils.Response{}}
}

func (h *AirsHandler) GetAirHome(c *gin.Context) {

	h.resp.Success(c, "Air IoT API Service")
}

func (h *AirsHandler) GetIndoorVal(c *gin.Context) {

	airList, err := h.airIot.AirThings()
	if err != nil {
		h.resp.BadRequest(c, err.Error())
		return
	}

	h.resp.Success(c, airList)
}
func (h *AirsHandler) GetIndoorValById(c *gin.Context) {

	filter := struct {
		DeviceSn  string `json:"device_sn" form:"device_sn"`
		StartDate string `json:"start_date" form:"start_date"`
		EndDate   string `json:"end_date" form:"end_date"`
	}{}

	//err := c.ShouldBindJSON(&filter)

	//err := c.ShouldBindUri(&filter)

	if c.ShouldBind(&filter) == nil {
		//log.Println(filter.DeviceSn)
		//log.Println(filter.StartDate)
		//log.Println(filter.EndDate)

		filters := &airs.Filter{
			DeviceSn: filter.DeviceSn,
			StartAt:  filter.StartDate,
			EndAt:    filter.EndDate,
		}

		result, err := h.airIot.AirThingsById(filter.DeviceSn, filters)
		if err != nil {
			h.resp.BadRequest(c, err.Error())
			return
		}

		h.resp.Success(c, result)
		//c.Header("Content-Type", "application/json")
		//
		////aorJson, _ := json.Marshal(result)
		////
		////c.AsciiJSON(200, string(aorJson))
		////c.AsciiJSON(200, string(aorJson))
		//
		//c.JSON(200, result)

		return
	}

	//if err != nil {
	//	h.resp.BadRequest(c, err.Error())
	//	return
	//}

	//log.Println(filter)

	c.JSON(400, filter)
	//deviceSn := c.Param("sn")

}
func (h *AirsHandler) GetAirToGrafana(c *gin.Context) {

	filter := struct {
		DeviceSn  string `json:"device_sn" form:"device_sn"`
		StartDate string `json:"start_date" form:"start_date"`
		EndDate   string `json:"end_date" form:"end_date"`
		Limit     int    `json:"limit" form:"limit"`
	}{}

	if c.ShouldBind(&filter) == nil {
		//log.Println(filter.DeviceSn)
		//log.Println(filter.StartDate)
		//log.Println(filter.EndDate)

		filters := &airs.Filter{
			DeviceSn: filter.DeviceSn,
			StartAt:  filter.StartDate,
			EndAt:    filter.EndDate,
			Limit:    filter.Limit,
		}

		result, err := h.airIot.AirThingsById2(filter.DeviceSn, filters)
		if err != nil {
			h.resp.BadRequest(c, err.Error())
			return
		}

		h.resp.Success(c, result)
		return
	}

	c.JSON(400, filter)

}
func (h *AirsHandler) GetAirToGrafanaMonitor(c *gin.Context) {

	filters := struct {
		DeviceSn  string `json:"device_sn" form:"device_sn"`
		StartDate string `json:"start_date" form:"start_date"`
		EndDate   string `json:"end_date" form:"end_date"`
		Mode      string `json:"mode" form:"mode"`
		Data      string `json:"data" form:"data"`
		Limit     int    `json:"limit" form:"limit"`
	}{}

	if c.ShouldBind(&filters) == nil {

		filter := &airs.Filter{
			DeviceSn: filters.DeviceSn,
			StartAt:  filters.StartDate,
			EndAt:    filters.EndDate,
			Mode:     filters.Mode,
			Data:     filters.Data,
			Limit:    filters.Limit,
		}

		regData := filters.Data
		switch regData {
		case "1000":

			rsRaw := []*airs.AirInGrafana{}
			rsDisplay := []*airs.AirReport{}

			var er error

			if filters.Mode == "raw" {
				rsRaw, er = h.airIot.AirThingsById2(filter.DeviceSn, filter)
				if er != nil {
					h.resp.BadRequest(c, er.Error())
					return
				}
				h.resp.Success(c, rsRaw)
				return
			} else {
				rsDisplay, er = h.airIot.AirThingsById(filter.DeviceSn, filter)
				if er != nil {
					h.resp.BadRequest(c, er.Error())
					return
				}
				h.resp.Success(c, rsDisplay)
				return
			}

		case "2000":

			//log.Println("reg 2000")

			result, err := h.airIot.DeviceThingsById(filter)
			if err != nil {
				h.resp.BadRequest(c, err)
				return
			}
			h.resp.Success(c, result)
			return

		case "3000":
			//log.Println("reg 3000")

			result, err := h.airIot.DeviceThingsOduVal(filter)
			if err != nil {
				h.resp.BadRequest(c, err)
				return
			}
			h.resp.Success(c, result)
			return

		case "4000":
			log.Println("reg 4000")
		default:
			log.Println("Default")
			h.resp.Success(c, filter)
			return
		}

	}
	h.resp.Success(c, filters)
	return
}
