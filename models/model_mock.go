package models

import (
	"github.com/jinzhu/gorm"
	mocket "github.com/selvatico/go-mocket"
)

// change the gorm client to use the mock db for testing
func InitMockModel() {
	mocket.Catcher.Register()
	db, _ = gorm.Open(mocket.DriverName, "connection_string")
}
