package handler

import (
	"Airnewnormal/domain/airs"
	"Airnewnormal/utils"
	"github.com/gin-gonic/gin"
)

type AirsHandler struct {
	airIot airs.AirService
	resp   utils.Response
}

func NewAirsHandler(airIot airs.AirService) AirsHandler {

	return AirsHandler{airIot: airIot, resp: utils.Response{}}
}

func (h *AirsHandler) GetIndoorVal(c *gin.Context) {

	airList, err := h.airIot.AirThings()
	if err != nil {
		h.resp.BadRequest(c, err.Error())
		return
	}
	
	h.resp.Success(c, airList)
}
