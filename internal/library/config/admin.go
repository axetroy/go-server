// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type admin struct {
	Domain string `json:"domain"` // HTTP API 绑定的域名
	Port   string `json:"port"`   // HTTP API 监听的端口
	Secret string `json:"secret"` // 管理员端密钥，用于加密/解密 token
}

var Admin admin

func init() {
	Admin.Port = dotenv.GetByDefault("ADMIN_HTTP_PORT", "9001")
	Admin.Domain = dotenv.GetByDefault("ADMIN_HTTP_DOMAIN", "localhost")
	Admin.Secret = dotenv.GetByDefault("ADMIN_TOKEN_SECRET_KEY", "admin")
}
