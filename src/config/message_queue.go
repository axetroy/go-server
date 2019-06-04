// Copyright 2019 Axetroy. All rights reserved. MIT license.
package config

import "os"

type messageQueue struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

var MessageQueue messageQueue

func init() {
	if MessageQueue.Host = os.Getenv("MSG_QUEUE_SERVER"); MessageQueue.Host == "" {
		MessageQueue.Host = "127.0.0.1"
	}
	if MessageQueue.Port = os.Getenv("MSG_QUEUE_PORT"); MessageQueue.Port == "" {
		MessageQueue.Port = "4150"
	}
}
