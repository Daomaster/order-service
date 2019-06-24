package main

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"order-service/api"
	"order-service/config"
	"order-service/models"
	"order-service/services/distance"
)

func init() {
	config.InitConfig()
	distance.InitGoogleMapCalculator()
	models.InitModel()
}

func main() {
	// init the router
	g := api.InitRouter()
	g.Use(gin.Logger())
	g.Use(gin.Recovery())

	// run on 8080 for the server
	err := g.Run(":8080")
	if err != nil {
		logrus.Fatal(err)
	}
}
