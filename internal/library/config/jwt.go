// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type jwt struct {
	Secret string `json:"secret"` // 用户端密钥，用于加密/解密 token
}

var Jwt jwt

func init() {
	Jwt.Secret = dotenv.GetByDefault("TOKEN_SECRET_KEY", "44JodlDOWk13f8a0&fKSDI*AP")
}
