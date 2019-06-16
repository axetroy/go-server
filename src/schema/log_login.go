// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

type LogLoginPure struct {
	Id      string `json:"id"`
	Uid     string `json:"uid"`
	Type    int    `json:"type"`
	Command int    `json:"command"`
	LastIp  string `json:"last_ip"`
	Client  string `json:"client"`
}

type LogLogin struct {
	LogLoginPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
