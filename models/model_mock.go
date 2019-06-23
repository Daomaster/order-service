package models

import (
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
)

func InitMockModel() {
	mocket.Catcher.Register()
	mocket.Catcher.Logging = true
	db, _ = gorm.Open(mocket.DriverName, "connection_string")
}
