package main

import (
	conf "Airnewnormal/config"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	server *gin.Engine
)

func init() {
	server = gin.Default()
}

func main() {
	config, _ := conf.LoadConfig(".")
	startGinServer(config)
}

func startGinServer(cf conf.Config) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{cf.Origin}
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"}
	server.Use(cors.New(corsConfig))
	server.Use(gin.Recovery())

	server.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"code":    "PAGE_NOT_FOUND",
			"message": "page not found",
		})
	})
	router := server.Group("/api/v1")

	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "OK"})
	})

	log.Fatal(server.Run(":8888"))

}
