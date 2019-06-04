// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import "os"

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
	SMTP.Host = os.Getenv("SMTP_SERVER")
	SMTP.Port = os.Getenv("SMTP_SERVER_PORT")
	SMTP.Username = os.Getenv("SMTP_USERNAME")
	SMTP.Password = os.Getenv("SMTP_PASSWORD")
	SMTP.Sender.Name = os.Getenv("SMTP_FROM_NAME")
	SMTP.Sender.Email = os.Getenv("SMTP_FROM_EMAIL")
}
