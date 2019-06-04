// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import "os"

type admin struct {
	Domain string `json:"domain"` // 管理员端 API 绑定的域名
	Port   string `json:"port"`   // 管理员端 API 监听的端口
	Secret string `json:"secret"` // 管理员端密钥，用于加密/解密 token
}

var Admin admin

func init() {
	if Admin.Port = os.Getenv("ADMIN_HTTP_PORT"); Admin.Port == "" {
		Admin.Port = "8081"
	}
	if Admin.Domain = os.Getenv("ADMIN_HTTP_DOMAIN"); Admin.Domain == "" {
		Admin.Domain = "http://127.0.0.1:" + Admin.Port
	}
	if Admin.Secret = os.Getenv("ADMIN_TOKEN_SECRET_KEY"); Admin.Secret == "" {
		Admin.Secret = "admin"
	}
}
