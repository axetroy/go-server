// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import (
	"github.com/axetroy/go-server/internal/service/dotenv"
)

type messageQueue struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var MessageQueue messageQueue

func init() {
	MessageQueue.Host = dotenv.GetByDefault("MSG_QUEUE_SERVER", "127.0.0.1")
	MessageQueue.Port = dotenv.GetByDefault("MSG_QUEUE_PORT", "4150")
}
