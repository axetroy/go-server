// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type resource struct {
	Domain string `json:"domain"` // HTTP 绑定的域名, 例如 https://example.com
	Port   string `json:"port"`   // HTTP API 监听的端口
}

var Resource resource

func init() {
	Resource.Port = dotenv.GetByDefault("RESOURCE_HTTP_PORT", "9004")
	Resource.Domain = dotenv.GetByDefault("RESOURCE_HTTP_DOMAIN", "localhost")
}
