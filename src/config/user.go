package config

import "os"

type user struct {
	Domain string `json:"domain"`
	Port   string `json:"port"`
	Secret string `json:"secret"`
}

var User user

func init() {
	if User.Port = os.Getenv("USER_HTTP_PORT"); User.Port == "" {
		User.Port = "8080"
	}
	if User.Domain = os.Getenv("USER_HTTP_DOMAIN"); User.Domain == "" {
		User.Domain = "http://127.0.0.1:" + User.Port
	}
	if User.Secret = os.Getenv("USER_TOKEN_SECRET_KEY"); User.Secret == "" {
		User.Secret = "user"
	}
}
