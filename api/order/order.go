package order

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
	"net/http"
	r "order-service/api/requests"
	"order-service/models"
	"order-service/pkgs/e"
	"strconv"
)

// handler for creating order
func CreateOrder(c *gin.Context) {
	var req r.CreateOrderRequest
	if err := c.BindJSON(&req); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, e.ErrOrderRequestInvalid)
		return
	}

	o, err := models.CreateOrder(req.Origin, req.Destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, e.CreateErr(e.ErrInternalError))
		return
	}

	c.JSON(http.StatusOK, o)
}

// handler for get list of the order
func GetOrders(c *gin.Context) {
	var req r.GetOrderRequest
	if err := c.BindQuery(&req); err != nil {
		logrus.Error(err)
		c.JSON(http.StatusBadRequest, e.CreateErr(e.ErrQueryStringInvalid))
		return
	}

	// make sure the req is valid
	if req.Page < 0 || req.Limit < 0 {
		c.JSON(http.StatusBadRequest, e.CreateErr(e.ErrQueryStringInvalid))
		return
	}

	os, err := models.GetOrders(req.Page, req.Limit)
	if err != nil && err != gorm.ErrRecordNotFound {
		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, e.CreateErr(e.ErrInternalError))
		return
	}

	c.JSON(http.StatusOK, os)
}

// handler for update an existing order
func UpdateOrder(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, e.ErrOrderRequestInvalid)
		return
	}

	var req r.TakeOrderRequest
	if err := c.BindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, e.ErrOrderRequestInvalid)
		return
	}

	// got anything other than TAKEN
	if req.Status != models.StatusTaken {
		c.JSON(http.StatusBadRequest, e.ErrOrderRequestInvalid)
		return
	}

	err = models.TakeOrder(id)
	if err != nil {
		// order is not found
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, e.CreateErr(e.ErrOrderNotExist))
			return
		}
		// order has already been taken
		if err == e.ErrOrderAlreadyTaken {
			c.JSON(http.StatusConflict, e.CreateErr(e.ErrOrderAlreadyTaken))
			return
		}

		logrus.Error(err)
		c.JSON(http.StatusInternalServerError, e.CreateErr(e.ErrInternalError))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": models.StatusSuccess,
	})
}
