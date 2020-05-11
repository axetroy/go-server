// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type wechat struct {
	AppID  string `json:"app_id"`
	Secret string `json:"secret"`
}

var Wechat wechat

func init() {
	Wechat.AppID = dotenv.Get("WECHAT_APP_ID")
	Wechat.Secret = dotenv.Get("WECHAT_SECRET")
}
