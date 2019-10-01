// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/src/service/dotenv"
)

type admin struct {
	Domain string `json:"domain"` // 管理员端 API 绑定的域名
	Port   string `json:"port"`   // 管理员端 API 监听的端口
	Secret string `json:"secret"` // 管理员端密钥，用于加密/解密 token
}

var Admin admin

func init() {
	Admin.Port = dotenv.GetByDefault("ADMIN_HTTP_PORT", "8081")
	Admin.Domain = dotenv.GetByDefault("ADMIN_HTTP_DOMAIN", "http://127.0.0.1:"+Admin.Port)
	Admin.Secret = dotenv.GetByDefault("ADMIN_TOKEN_SECRET_KEY", "admin")
}
