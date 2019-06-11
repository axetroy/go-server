// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message_queue

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/config"
	"github.com/axetroy/go-server/service/email"
	"github.com/axetroy/go-server/service/redis"
	"github.com/nsqio/go-nsq"
	"sync"
	"time"
)

type Topic string
type Chanel string

var (
	TopicSendEmail  Topic       = "send_email"
	ChanelSendEmail Chanel      = "send_email"
	Address         string      // 消息队列地址
	Config          *nsq.Config // 消息队列的配置
)

type SendActivationEmailBody struct {
	Email string `json:"email"` // 要发送的邮箱
	Code  string `json:"code"`  // 发送的激活码
}

func init() {
	host := config.MessageQueue.Host
	port := config.MessageQueue.Port

	Address = host + ":" + port

	Config = nsq.NewConfig()
	Config.DialTimeout = time.Second * 5
	Config.MsgTimeout = time.Second * 10
	Config.ReadTimeout = time.Second * 15
	Config.WriteTimeout = time.Second * 10
	Config.HeartbeatInterval = time.Second * 10
}

func RunMessageQueueConsumer() {
	var (
		err error
	)

	wg := &sync.WaitGroup{}

	wg.Add(1)

	_, err = CreateConsumer(TopicSendEmail, ChanelSendEmail, nsq.HandlerFunc(func(message *nsq.Message) error {

		body := SendActivationEmailBody{}

		if err := json.Unmarshal(message.Body, &body); err != nil {
			return err
		}

		mailer := email.NewMailer()

		// 发送邮件
		if err := mailer.SendActivationEmail(body.Email, body.Code); err != nil {
			// 邮件没发出去的话，删除 redis 的 key
			_ = redis.ActivationCodeClient.Del(body.Code).Err()
		}

		fmt.Printf("发送验证码 %s 到 %s\n", body.Code, body.Email)

		return nil
	}))

	if err != nil {
		panic(err)
	}

	wg.Wait()

	return
}
