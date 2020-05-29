// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type user struct {
	Secret string `json:"secret"` // 用户端密钥，用于加密/解密 token
}

var User user

func init() {
	User.Secret = dotenv.GetByDefault("USER_TOKEN_SECRET_KEY", "user")
}
