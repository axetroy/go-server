// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type user struct {
	Domain string `json:"domain"` // HTTP 绑定的域名, 例如 https://example.com
	Port   string `json:"port"`   // HTTP 监听的端口
	Secret string `json:"secret"` // 用户端密钥，用于加密/解密 token
}

var User user

func init() {
	User.Port = dotenv.GetByDefault("USER_HTTP_PORT", "9000")
	User.Domain = dotenv.GetByDefault("USER_HTTP_DOMAIN", "localhost")
	User.Secret = dotenv.GetByDefault("USER_TOKEN_SECRET_KEY", "user")
}
