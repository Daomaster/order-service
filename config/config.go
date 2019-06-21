package config

import (
	"fmt"
	"github.com/spf13/viper"
)

var config *Configuration

const ConnectionStringFormat = "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local"

type Configuration struct {
	MapConfig *MapConfiguration
	DbConfig *DbConfiguration
}

type MapConfiguration struct {
	apiKey string
}

// return the map api key
func (m MapConfiguration) GetMapApiKey() string {
	return m.apiKey
}

type DbConfiguration struct {
	hostname   string
	port       int
	username   string
	password   string
	schemaName string
}

// return the whole connection string
func (d DbConfiguration) GetConnectionString() string {
	return fmt.Sprintf(ConnectionStringFormat, d.username, d.password, d.hostname, d.port, d.schemaName)
}

// getter of the config var
func GetConfig() *Configuration {
	return config
}

// initialization function for the config
func InitConfig() {
	v := viper.New()
	v.AddConfigPath(".")
	v.SetConfigFile("config.json")
	v.SetConfigType("json")

	// default values
	v.SetDefault("MYSQL_SCHEMA", "order-service")
	v.SetDefault("MYSQL_PORT", 3306)

	err := v.ReadInConfig()

	// local config file does not exist
	if err != nil {
		v.BindEnv("MYSQL_ROOT_PWD")
		v.BindEnv("MYSQL_HOSTNAME")
		v.BindEnv("MYSQL_USER")
		v.BindEnv("MAP_API_KEY")
	} else {
		// overwrite if env is present
		v.AutomaticEnv()
	}

	config = new(Configuration)

	var dbConfig DbConfiguration
	dbConfig.port = v.GetInt("MYSQL_PORT")
	dbConfig.hostname = v.GetString("MYSQL_HOSTNAME")
	dbConfig.username = v.GetString("MYSQL_USER")
	dbConfig.password = v.GetString("MYSQL_ROOT_PWD")
	dbConfig.schemaName = v.GetString("MYSQL_SCHEMA")

	var mapConfig MapConfiguration
	mapConfig.apiKey = v.GetString("MAP_API_KEY")

	config.DbConfig = &dbConfig
	config.MapConfig = &mapConfig
}
