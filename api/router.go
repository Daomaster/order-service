package api

import (
	"github.com/gin-gonic/gin"
	"order-service/api/order"
)

// function for initialize the routes for gin
func InitRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	orderRoute := r.Group(`/orders`)

	orderRoute.Use()
	{
		// get orders
		orderRoute.GET("/", order.GetOrders)

		// update status of an order
		orderRoute.PATCH("/:id", order.UpdateOrder)

		// create a new order
		orderRoute.POST("/", order.CreateOrder)
	}

	return r
}