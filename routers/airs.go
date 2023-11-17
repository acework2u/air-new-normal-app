package routers

import (
	"Airnewnormal/handler"
	"github.com/gin-gonic/gin"
)

type AirsRouter struct {
	airHandler handler.AirsHandler
}

func NewAirsRouter(airsHandler handler.AirsHandler) AirsRouter {
	return AirsRouter{airHandler: airsHandler}
}

func (r *AirsRouter) AirRoute(rg *gin.RouterGroup) {
	router := rg.Group("air-iot")
	router.GET("/", r.airHandler.GetIndoorVal)

}
