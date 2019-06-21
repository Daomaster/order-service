package order

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
	r "order-service/api/requests"
	res_error "order-service/api/res-error"
	"order-service/models"
)

// handler for creating order
func CreateOrder(c *gin.Context) {
	var req r.CreateOrderRequest
	if err := c.BindJSON(&req); err != nil {
		logrus.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	o, err := models.CreateOrder(req.Origin, req.Destination)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, o)
}

// handler for get list of the order
func GetOrders(c *gin.Context) {
	var req r.GetOrderRequest
	if err := c.BindQuery(&req); err != nil {
		logrus.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	if req.Page < 0 || req.Limit < 0 {
		c.Status(http.StatusBadRequest)
		return
	}

	os, err := models.GetOrders(req.Page, req.Limit)
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Error(err)
		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, os)
}

// handler for update an existing order
func UpdateOrder(c *gin.Context) {
	var req r.TakeOrderRequest
	if err := c.BindUri(&req); err != nil {
		logrus.Error(err)
		c.Status(http.StatusBadRequest)
		return
	}

	err := models.TakeOrder(req.ID)
	if err != nil {
		logrus.Error(err)

		if err == gorm.ErrRecordNotFound {
			c.Status(http.StatusNotFound)
			return
		}

		if err == res_error.ErrOrderAlreadyTaken {
			c.Status(http.StatusConflict)
			return
		}

		c.Status(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": models.StatusSuccess,
	})
}
