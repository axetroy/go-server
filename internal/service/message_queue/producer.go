// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message_queue

import (
	"encoding/json"
	"errors"
	"github.com/axetroy/go-server/internal/service/notify"
	"github.com/nsqio/go-nsq"
	"time"
)

var (
	producer *nsq.Producer
)

func init() {
	if err := createProducer(); err != nil {
		panic(err)
	}
}

// 初始化生产者
func createProducer() (err error) {
	if producer, err = nsq.NewProducer(Address, Config); err != nil {
		return
	}
	return
}

func DeferredPublish(topic Topic, delay time.Duration, message []byte) (err error) {
	var (
		maxConnectTimes = 5
		connectTimes    = 0
	)

	// 确保链接可用
	for {
		if producer.Ping() == nil {
			break
		}
		if connectTimes >= maxConnectTimes {
			err = errors.New("publish timeout")
			return
		}
		connectTimes = connectTimes + 1
	}

	//不能发布空串，否则会导致 error
	if len(message) == 0 {
		err = errors.New("message can not be empty")
		return
	}

	if err = producer.DeferredPublish(string(topic), delay, message); err != nil {
		return
	}

	return
}

// 发布消息
func Publish(topic Topic, message []byte) (err error) {
	var (
		maxConnectTimes = 5
		connectTimes    = 0
	)

	// 确保链接可用
	for {
		if producer.Ping() == nil {
			break
		}
		if connectTimes >= maxConnectTimes {
			err = errors.New("publish timeout")
			return
		}
		connectTimes = connectTimes + 1
	}

	//不能发布空串，否则会导致 error
	if len(message) == 0 {
		err = errors.New("message can not be empty")
		return
	}

	if err = producer.Publish(string(topic), message); err != nil {
		return
	}

	return
}

type PayloadPublishSystemNotification struct {
	NotificationID string `json:"notification_id" valid:"required"`
}

type PayloadPublishCheckUserLoginStatus struct {
	UserID string `json:"user_id" valid:"required"`
}

type PayloadToAllUsers struct {
	Title   string                 `json:"title" valid:"required"`   // 推送的标题
	Content string                 `json:"content" valid:"required"` // 推送的内容
	Data    map[string]interface{} `json:"data"`                     // 附带给 APP 的数据
}

type PayloadToSpecificUsers struct {
	UserID  []string               `json:"user_id" valid:"required"` // 要指定的推送用户 ID
	Title   string                 `json:"title" valid:"required"`   // 推送的标题
	Content string                 `json:"content" valid:"required"` // 推送的内容
	Data    map[string]interface{} `json:"data"`                     // 附带给 APP 的数据
}

// 推送 - 系统通知
func PublishSystemNotify(notificationID string) error {
	body := BodySendNotify{
		Event: notify.EventSendNotifyToUserNewNotification,
		Payload: PayloadPublishSystemNotification{
			NotificationID: notificationID,
		},
	}

	b, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return DeferredPublish(TopicPushNotify, 0, b)
}

// 推送 - 检查用户的登录状态
func PublishCheckUserLogin(userID string) error {
	body := BodySendNotify{
		Event: notify.EventSendNotifyCheckUserLoginStatus,
		Payload: PayloadPublishCheckUserLoginStatus{
			UserID: userID,
		},
	}

	b, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return DeferredPublish(TopicPushNotify, time.Second*60, b)
}

// 发送到消息队列 - 发送推送给所有用户
func PublishNotifyToAllUser(title string, content string, delay time.Duration, data map[string]interface{}) error {
	body := BodySendNotify{
		Event: notify.EventSendNotifyToAllUser,
		Payload: PayloadToAllUsers{
			Title:   title,
			Content: content,
			Data:    data,
		},
	}

	b, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return DeferredPublish(TopicPushNotify, delay, b)
}

// 发送到消息队列 - 发送推送给特定用户
func PublishNotifyToSpecificUser(userId []string, title string, content string, delay time.Duration, data map[string]interface{}) error {
	body := BodySendNotify{
		Event: notify.EventSendNotifyToCustomUser,
		Payload: PayloadToSpecificUsers{
			UserID:  userId,
			Title:   title,
			Content: content,
			Data:    data,
		},
	}

	b, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return DeferredPublish(TopicPushNotify, delay, b)
}
