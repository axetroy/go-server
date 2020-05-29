// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type admin struct {
	Secret string `json:"secret"` // 管理员端密钥，用于加密/解密 token
}

var Admin admin

func init() {
	Admin.Secret = dotenv.GetByDefault("ADMIN_TOKEN_SECRET_KEY", "admin")
}
