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
	router.GET("", r.airHandler.GetAirHome)
	router.GET("/grafana", r.airHandler.GetAirToGrafana)
	router.GET("/indoor", r.airHandler.GetIndoorValById)
	router.GET("/indoor/:id", r.airHandler.GetIndoorValById)
	router.GET("/indoors", r.airHandler.GetIndoorVal)

	router2 := rg.Group("grafana-air-service")
	router2.GET("", r.airHandler.GetAirToGrafanaMonitor)
}
