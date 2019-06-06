// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/src/service/dotenv"
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
	if Database.Driver = dotenv.Get("DB_DRIVER"); Database.Driver == "" {
		Database.Driver = "postgres"
	}
	if Database.Host = dotenv.Get("DB_HOST"); Database.Host == "" {
		Database.Host = "localhost"
	}
	if Database.Port = dotenv.Get("DB_PORT"); Database.Port == "" {
		Database.Port = "65432"
	}
	if Database.DatabaseName = dotenv.Get("DB_NAME"); Database.DatabaseName == "" {
		Database.DatabaseName = "gotest"
	}
	if Database.Username = dotenv.Get("DB_USERNAME"); Database.Username == "" {
		Database.Username = "gotest"
	}
	if Database.Password = dotenv.Get("DB_PASSWORD"); Database.Password == "" {
		Database.Password = "gotest"
	}
	if Database.Sync = dotenv.Get("DB_SYNC"); Database.Sync == "" {
		Database.Sync = "on"
	}
}
