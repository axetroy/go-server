// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type sender struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type smtp struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Sender   sender `json:"sender"`
}

var SMTP smtp

func init() {
	SMTP.Host = dotenv.Get("SMTP_SERVER")
	SMTP.Port = dotenv.Get("SMTP_SERVER_PORT")
	SMTP.Username = dotenv.Get("SMTP_USERNAME")
	SMTP.Password = dotenv.Get("SMTP_PASSWORD")
	SMTP.Sender.Name = dotenv.Get("SMTP_FROM_NAME")
	SMTP.Sender.Email = dotenv.Get("SMTP_FROM_EMAIL")
}
