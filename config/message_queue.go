// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/service/dotenv"
)

type messageQueue struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var MessageQueue messageQueue

func init() {
	if MessageQueue.Host = dotenv.Get("MSG_QUEUE_SERVER"); MessageQueue.Host == "" {
		MessageQueue.Host = "127.0.0.1"
	}
	if MessageQueue.Port = dotenv.Get("MSG_QUEUE_PORT"); MessageQueue.Port == "" {
		MessageQueue.Port = "4150"
	}
}
