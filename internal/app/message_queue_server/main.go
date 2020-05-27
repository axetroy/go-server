// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message_queue_server

import (
	"github.com/axetroy/go-server/internal/app/message_queue_server/handler"
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/nsqio/go-nsq"
	"sync"
)

func CreateConsumer(handler handler.Handler) (*nsq.Consumer, error) {
	return message_queue.CreateConsumer(handler.GetTopic(), handler.GetChannel(), nsq.HandlerFunc(func(message *nsq.Message) error {
		return handler.OnMessage(message)
	}))
}

func RunMessageQueueConsumer() ([]*nsq.Consumer, error) {
	consumers := make([]*nsq.Consumer, 0)
	wg := &sync.WaitGroup{}

	wg.Add(100)

	emailConsumer, err := CreateConsumer(handler.NewEmailHandler(message_queue.TopicSendEmail, message_queue.ChanelSendEmail))

	if err != nil {
		return consumers, err
	} else {
		consumers = append(consumers, emailConsumer)
	}

	notifyConsumer, err := CreateConsumer(handler.NewNotifyHandler(message_queue.TopicPushNotify, message_queue.ChanelPushNotify))

	if err != nil {
		return consumers, err
	} else {
		consumers = append(consumers, notifyConsumer)
	}

	wg.Wait()

	return consumers, nil
}
