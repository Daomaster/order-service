package distance

import (
	"context"
	"github.com/sirupsen/logrus"
	"googlemaps.github.io/maps"
	"order-service/config"
	"order-service/pkgs/e"
	"strings"
)

type googleMapCalculator struct {
	client *maps.Client
}

// initialize the google map client
func InitGoogleMapCalculator() {
	apiKey := config.GetConfig().MapConfig.GetMapApiKey()

	// init google map api client
	c, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		logrus.Fatal(err)
	}

	var googleCalc googleMapCalculator
	googleCalc.client = c

	calc = &googleCalc
}

// calculate the distance between
func (c *googleMapCalculator) Calculate(src []string, des []string) (int, error) {
	srcStr := strings.Join(src, ",")
	desStr := strings.Join(des, ",")

	req := new(maps.DistanceMatrixRequest)
	req.Origins = append(req.Origins, srcStr)
	req.Destinations = append(req.Destinations, desStr)

	// use the distance matrix api
	res, err := c.client.DistanceMatrix(context.Background(), req)
	if err != nil {
		return 0, err
	}

	// use the first result since there is no specific info provided
	el := res.Rows[0].Elements[0]

	// if google can't find a route
	if el.Status != "OK" {
		return 0, e.ErrDistanceUnknown
	}

	// use the the first result
	result := el.Distance.Meters

	return result, nil
}
