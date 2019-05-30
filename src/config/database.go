package config

import (
	"os"
)

type database struct {
	Host         string `json:"host"`
	Port         string `json:"port"`
	Driver       string `json:"driver"`
	DatabaseName string `json:"database_name"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Sync         string `json:"sync"`
}

var Database database

func init() {
	if Database.Driver = os.Getenv("DB_DRIVER"); Database.Driver == "" {
		Database.Driver = "postgres"
	}
	if Database.Host = os.Getenv("DB_HOST"); Database.Host == "" {
		Database.Host = "localhost"
	}
	if Database.Port = os.Getenv("DB_PORT"); Database.Port == "" {
		Database.Port = "65432"
	}
	if Database.DatabaseName = os.Getenv("DB_NAME"); Database.DatabaseName == "" {
		Database.DatabaseName = "gotest"
	}
	if Database.Username = os.Getenv("DB_USERNAME"); Database.Username == "" {
		Database.Username = "gotest"
	}
	if Database.Password = os.Getenv("DB_PASSWORD"); Database.Password == "" {
		Database.Password = "gotest"
	}
	if Database.Sync = os.Getenv("DB_SYNC"); Database.Sync == "" {
		Database.Sync = "on"
	}
}
