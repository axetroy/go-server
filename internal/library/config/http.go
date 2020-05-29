// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

var Http http

type http struct {
	Host   string `json:"host"`   // HTTP 监听的地址
	Port   string `json:"port"`   // HTTP 监听的端口
	Domain string `json:"domain"` // HTTP 绑定的域名, 例如 https://example.com
}

func init() {
	Http.Host = dotenv.GetByDefault("HOST", "127.0.0.1")
	Http.Port = dotenv.GetByDefault("PORT", "80")
	Http.Domain = dotenv.GetByDefault("DOMAIN", "https://example.com")
}
