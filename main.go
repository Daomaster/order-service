package main

import (
	"github.com/sirupsen/logrus"
	"order-service/api"
	"order-service/config"
	"order-service/models"
)

func init()  {
	config.InitConfig()

	err := models.Init()
	if err != nil {
		logrus.Fatal(err)
	}
}

func main() {
	// init the router
	g := api.InitRouter()

	// run on 8000 for the server
	err := g.Run(":8000")
	if err != nil {
		logrus.Fatal(err)
	}
}
