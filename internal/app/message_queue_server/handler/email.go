// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package handler

import (
	"context"
	"encoding/json"
	"github.com/axetroy/go-server/internal/service/email"
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/nsqio/go-nsq"
	"log"
)

type EmailHandler struct {
	topic  message_queue.Topic
	chanel message_queue.Chanel
}

func NewEmailHandler(topic message_queue.Topic, chanel message_queue.Chanel) *NotifyHandler {
	return &NotifyHandler{
		topic:  topic,
		chanel: chanel,
	}
}

func (h *EmailHandler) GetTopic() message_queue.Topic {
	return h.topic
}

func (h *EmailHandler) GetChannel() message_queue.Chanel {
	return h.chanel
}

func (h *EmailHandler) OnMessage(message *nsq.Message) error {
	body := message_queue.BodySendActivationEmail{}

	if err := json.Unmarshal(message.Body, &body); err != nil {
		return err
	}

	mailer, err := email.NewMailer()

	if err != nil {
		return err
	}

	// 发送邮件
	if err := mailer.SendActivationEmail(body.Email, body.Code); err != nil {
		// 邮件没发出去的话，删除 redis 的 key
		_ = redis.ClientActivationCode.Del(context.Background(), body.Code).Err()
	}

	log.Printf("发送验证码 %s 到 %s\n", body.Code, body.Email)

	return nil
}
