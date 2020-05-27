// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message_queue

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/service/notify"
	"github.com/nsqio/go-nsq"
	"net"
	"time"
)

type Topic string
type Chanel string

var (
	TopicSendEmail   Topic       = "topic_send_email"
	ChanelSendEmail  Chanel      = "chanel_send_email"
	TopicPushNotify  Topic       = "topic_push_notify"
	ChanelPushNotify Chanel      = "chanel_push_notify"
	Address                      = net.JoinHostPort(config.MessageQueue.Host, config.MessageQueue.Port) // 消息队列地址
	Config           *nsq.Config                                                                        // 消息队列的配置
)

type BodySendActivationEmail struct {
	Email string `json:"email" valid:"required"` // 要发送的邮箱
	Code  string `json:"code" valid:"required"`  // 发送的激活码
}

type BodySendNotify struct {
	Event   notify.Event `json:"event" valid:"required"`   // 事件名称
	Payload interface{}  `json:"payload" valid:"required"` // 数据体
}

func (c *BodySendNotify) ToByte() ([]byte, error) {
	b, err := json.Marshal(c)

	if err != nil {
		return nil, err
	}

	return b, nil
}

func init() {
	Config = nsq.NewConfig()
	Config.DialTimeout = time.Second * 60
	Config.MsgTimeout = time.Second * 60
	Config.ReadTimeout = time.Second * 60
	Config.WriteTimeout = time.Second * 60
	Config.HeartbeatInterval = time.Second * 10
}
