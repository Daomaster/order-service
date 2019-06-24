package models

import (
	"github.com/jinzhu/gorm"
	"order-service/pkgs/e"
	"order-service/services/distance"
)

const (
	StatusUnassigned = "UNASSIGNED"
	StatusTaken      = "TAKEN"
	StatusSuccess    = "SUCCESS"
)

type Order struct {
	ID       int64  `gorm:"PRIMARY_KEY;AUTO_INCREMENT" json:"id"`
	Distance int    `json:"distance"`
	Status   string `json:"status"`
}

// function to create an order base on the src to des
func CreateOrder(src []string, des []string) (*Order, error) {
	calc := distance.GetCalculator()

	// calculate the distance
	d, err := calc.Calculate(src, des)
	if err != nil {
		return nil, err
	}

	// create the order in the db
	o := Order{Distance: d, Status: StatusUnassigned}
	if err := db.Create(&o).Error; err != nil {
		return nil, err
	}

	return &o, nil
}

// function to take order based on the id provided
func TakeOrder(id int64) error {
	var o Order

	// check if there is a order based on the id
	err := db.Where("id = ?", id).First(&o).Error
	if err != nil {
		return err
	}

	// check if it is taken
	if o.Status == StatusTaken {
		return e.ErrOrderAlreadyTaken
	}

	// found and modify status
	if err := db.Model(&o).Update("status", StatusTaken).Error; err != nil {
		return err
	}

	return nil
}

// function to retrieve paged orders
func GetOrders(page int, limit int) ([]*Order, error) {
	var orders []*Order

	err := db.Offset(page * limit).Limit(limit).Find(&orders).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return orders, nil
}
