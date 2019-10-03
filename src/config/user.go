// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/src/service/dotenv"
)

type TLS struct {
	Cert string `json:"cert"` // 证书文件
	Key  string `json:"key"`  // Key 文件
}

type user struct {
	Domain string `json:"domain"` // 用户端 API 绑定的域名, 例如 https://example.com
	Port   string `json:"port"`   // 用户端 API 监听的端口
	Secret string `json:"secret"` // 用户端密钥，用于加密/解密 token
	TLS    *TLS   `json:"tls"`
}

var User user

func init() {
	User.Port = dotenv.GetByDefault("USER_HTTP_PORT", "8080")
	User.Domain = dotenv.GetByDefault("USER_HTTP_DOMAIN", "localhost")
	User.Secret = dotenv.GetByDefault("USER_TOKEN_SECRET_KEY", "user")

	TlsCert := dotenv.GetByDefault("USER_TLS_CERT", "")
	TlsKey := dotenv.GetByDefault("USER_TLS_KEY", "")

	if TlsCert != "" && TlsKey != "" {
		User.TLS = &TLS{
			Cert: TlsCert,
			Key:  TlsKey,
		}
	}
}
