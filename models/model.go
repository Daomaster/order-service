package models

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"order-service/config"
)

var client *Client

type Client struct {
	db *sql.DB
}

// initialize the client with mysql connection
func Init() error {
	cString := config.GetConfig().DbConfig.GetConnectionString()

	db, err := sql.Open("mysql", cString)
	if err != nil {
		return err
	}

	c := new(Client)
	c.db = db

	client = c

	return nil
}