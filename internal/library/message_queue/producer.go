// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message_queue

import (
	"encoding/json"
	"errors"
	"github.com/axetroy/go-server/internal/schema"
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

// 发送到消息队列 - 用户登录异常
func PublishNotifyWhenLoginAbnormal(userInfo schema.ProfilePublic, delay time.Duration) error {
	body := SendNotifyBody{
		Event:   notify.SendNotifyEventSendNotifyToLoginAbnormalUser,
		Payload: userInfo,
	}

	b, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return DeferredPublish(TopicPushNotify, delay, b)
}

// 发送到消息队列 - 发送推送给所有用户
func PublishNotifyToAllUser(title string, content string, delay time.Duration) error {
	body := SendNotifyBody{
		Event: notify.SendNotifyEventSendNotifyToAllUser,
		Payload: NotifyPayloadToAllUsers{
			Title:   title,
			Content: content,
		},
	}

	b, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return DeferredPublish(TopicPushNotify, delay, b)
}

// 发送到消息队列 - 发送推送给特定用户
func PublishNotifyToSpecificUser(userId []string, title string, content string, delay time.Duration) error {
	body := SendNotifyBody{
		Event: notify.SendNotifyEventSendNotifyToAllUser,
		Payload: NotifyPayloadToSpecificUsers{
			UserID:  userId,
			Title:   title,
			Content: content,
		},
	}

	b, err := json.Marshal(body)

	if err != nil {
		return err
	}

	return DeferredPublish(TopicPushNotify, delay, b)
}
