// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package handler

import (
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/nsqio/go-nsq"
)

type Handler interface {
	GetTopic() message_queue.Topic
	GetChannel() message_queue.Chanel
	OnMessage(message *nsq.Message) error
}
