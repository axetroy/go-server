// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type redis struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Password string `json:"password"`
}

var Redis redis

func init() {
	Redis.Host = dotenv.GetByDefault("REDIS_SERVER", "127.0.0.1")
	Redis.Port = dotenv.GetByDefault("REDIS_PORT", "6379")
	Redis.Password = dotenv.GetByDefault("REDIS_PASSWORD", "password")
}
