// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package handler

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/axetroy/go-server/internal/service/notify"
	"github.com/nsqio/go-nsq"
	"log"
)

type NotifyHandler struct {
	topic  message_queue.Topic
	chanel message_queue.Chanel
}

func NewNotifyHandler(topic message_queue.Topic, chanel message_queue.Chanel) *NotifyHandler {
	return &NotifyHandler{
		topic:  topic,
		chanel: chanel,
	}
}

func (h *NotifyHandler) GetTopic() message_queue.Topic {
	return h.topic
}

func (h *NotifyHandler) GetChannel() message_queue.Chanel {
	return h.chanel
}

// 发送给所有用户
func (h *NotifyHandler) handlerSendToAllUser(payload interface{}) error {
	b, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	var data message_queue.PayloadToAllUsers

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if err := validator.ValidateStruct(data); err != nil {
		return err
	}

	if err := notify.Notify.SendNotifyToAllUser(data.Title, data.Content, data.Data); err != nil {
		return err
	}

	return nil
}

// 发送给指定用户
func (h *NotifyHandler) handlerSendToCustomUser(payload interface{}) error {
	b, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	var data message_queue.PayloadToSpecificUsers

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if err := validator.ValidateStruct(data); err != nil {
		return err
	}

	if err := notify.Notify.SendNotifyToCustomUser(data.UserID, data.Title, data.Content, data.Data); err != nil {
		return err
	}

	return nil
}

// 发送给指定用户
func (h *NotifyHandler) handlerCheckUserLoginStatus(payload interface{}) error {
	b, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	var data message_queue.PayloadPublishCheckUserLoginStatus

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if err := validator.ValidateStruct(data); err != nil {
		return err
	}

	// 发送推送给用户
	if err := notify.Notify.SendNotifyToUserForLoginStatus(data.UserID); err != nil {
		return err
	}
	return nil
}

// 推送通知 - 用户有新的系统通知
func (h *NotifyHandler) handlerSendNewSystemNotification(payload interface{}) error {
	b, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	var data message_queue.PayloadPublishSystemNotification

	if err := json.Unmarshal(b, &data); err != nil {
		return err
	}

	if err := validator.ValidateStruct(data); err != nil {
		return err
	}

	// 发送推送给用户
	if err := notify.Notify.SendNotifySystemNotificationToUser(data.NotificationID); err != nil {
		return err
	}
	return nil
}

func (h *NotifyHandler) OnMessage(message *nsq.Message) error {
	body := message_queue.BodySendNotify{}

	if err := json.Unmarshal(message.Body, &body); err != nil {
		return err
	}

	log.Println(string(message.Body))

	switch body.Event {
	// 推送通知 - 用户有新的系统通知
	case notify.EventSendNotifyToUserNewNotification:
		return h.handlerSendNewSystemNotification(body.Payload)
	// 推送一个自定义通知给所有用户
	case notify.EventSendNotifyToAllUser:
		return h.handlerSendToAllUser(body.Payload)
	// 推送一个自定义通知给指定用户
	case notify.EventSendNotifyToCustomUser:
		return h.handlerSendToCustomUser(body.Payload)
	// 发送给指定用户，登录异常
	case notify.EventSendNotifyCheckUserLoginStatus:
		return h.handlerCheckUserLoginStatus(body.Payload)
	default:
		return nil
	}
}
