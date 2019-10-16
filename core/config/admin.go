// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/core/service/dotenv"
)

type admin struct {
	Domain string `json:"domain"` // 管理员端 API 绑定的域名
	Port   string `json:"port"`   // 管理员端 API 监听的端口
	Secret string `json:"secret"` // 管理员端密钥，用于加密/解密 token
	TLS    *TLS   `json:"tls"`
}

var Admin admin

func init() {
	Admin.Port = dotenv.GetByDefault("ADMIN_HTTP_PORT", "8081")
	Admin.Domain = dotenv.GetByDefault("ADMIN_HTTP_DOMAIN", "localhost")
	Admin.Secret = dotenv.GetByDefault("ADMIN_TOKEN_SECRET_KEY", "admin")

	TlsCert := dotenv.GetByDefault("ADMIN_TLS_CERT", "")
	TlsKey := dotenv.GetByDefault("ADMIN_TLS_KEY", "")

	if TlsCert != "" && TlsKey != "" {
		User.TLS = &TLS{
			Cert: TlsCert,
			Key:  TlsKey,
		}
	}
}
