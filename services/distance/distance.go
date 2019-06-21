package distance

import (
	"context"
	"github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
	"order-service/config"
	"strings"
)

var calculator *Calculator

type Calculator struct {
	client *maps.Client
}

// initialize the google map client
func InitCalculator() {
	apiKey := config.GetConfig().MapConfig.GetMapApiKey()

	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		logrus.Fatal(err)
	}

	calculator = new(Calculator)
	calculator.client = c
}

// calculate the distance between
func Calculate(src []string, des []string) (int, error) {
	srcStr := strings.Join(src, ",")
	desStr := strings.Join(des, ",")

	req := new(maps.DistanceMatrixRequest)
	req.Origins = append(req.Origins, srcStr)
	req.Destinations = append(req.Destinations, desStr)

	res, err := calculator.client.DistanceMatrix(context.Background(), req)
	if err != nil {
		return 0, err
	}

	// use the the first result
	result := res.Rows[0].Elements[0].Distance.Meters

	return result, nil
}
