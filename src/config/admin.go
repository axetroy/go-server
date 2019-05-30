package config

import "os"

type admin struct {
	Domain string `json:"domain"`
	Port   string `json:"port"`
	Secret string `json:"secret"`
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
