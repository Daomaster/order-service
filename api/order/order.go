package order

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	r "order-service/api/requests"
)

// handler for creating order
func CreateOrder(c *gin.Context) {
	var req r.Order
	if c.BindJSON(&req) != nil {
		c.Status(http.StatusBadRequest)
	}

	logrus.Info(req)
}

// handler for get list of the order
func GetOrders(c *gin.Context) {
	page := c.Query("page")
	limit := c.Query("limit")

	if page == "" || limit == "" {
		c.Status(http.StatusBadRequest)
	}

	logrus.Info(page, limit)
}

// handler for update an existing order
func UpdateOrder(c *gin.Context) {
	orderID := c.Param("id")
	logrus.Info(orderID)
}
