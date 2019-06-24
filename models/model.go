package models

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"order-service/config"
)

var db *gorm.DB

// initialize the client with mysql connection
func InitModel() {
	cString := config.GetConfig().DbConfig.GetConnectionString()

	var err error
	db, err = gorm.Open("mysql", cString)
	if err != nil {
		logrus.Fatal(err)
	}

	// migration on all the tables
	migrate()
}

// initialize the tables based on the model if not exist
func migrate() {
	db.AutoMigrate(&Order{})
}
