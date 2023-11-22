package main

import (
	conf "Airnewnormal/config"
	"Airnewnormal/domain/airs"
	"Airnewnormal/handler"
	"Airnewnormal/repository"
	"Airnewnormal/routers"
	"context"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
)

var (
	server               *gin.Engine
	airsThingsCollection *mongo.Collection
	ctx                  context.Context
	mongoClient          *mongo.Client

	// Airs
	airsHandler handler.AirsHandler
	AirsRouter  routers.AirsRouter
)

func init() {

	//load Env
	envConf, _ := conf.LoadConfig(".")
	//DB Connected
	mongoClient = conf.ConnectDB(envConf.DBUrl)

	//Airs
	airsThingsCollection = conf.GetCollection(mongoClient, "airs_things")
	airsThingsRepo := repository.NewAirRepositoryDB(ctx, airsThingsCollection)

	airService := airs.NewAirThingsService(airsThingsRepo)
	airsHandler = handler.NewAirsHandler(airService)
	AirsRouter = routers.NewAirsRouter(airsHandler)

	ctx = context.TODO()
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
	corsConfig.AllowHeaders = []string{"Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With,grafana backend/server,X-Grafana-Org-Id"}
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

	//AirThings
	AirsRouter.AirRoute(router)

	log.Fatal(server.Run(":8888"))

}
