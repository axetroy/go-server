// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message_queue

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/email"
	"github.com/axetroy/go-server/internal/service/notify"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/nsqio/go-nsq"
	"log"
	"net"
	"sync"
	"time"
)

type Topic string
type Chanel string

var (
	TopicSendEmail   Topic       = "send_email"
	ChanelSendEmail  Chanel      = "send_email"
	TopicPushNotify  Topic       = "push_notify"
	ChanelPushNotify Chanel      = "push_notify"
	Address          string      // 消息队列地址
	Config           *nsq.Config // 消息队列的配置
)

type SendActivationEmailBody struct {
	Email string `json:"email" valid:"required~请输入邮箱"` // 要发送的邮箱
	Code  string `json:"code" valid:"required~请输入激活码"` // 发送的激活码
}

type SendNotifyBody struct {
	Event   notify.SendNotifyEvent `json:"event" valid:"required~请输入事件"`    // 事件名称
	Payload interface{}            `json:"payload" valid:"required~请输入数据体"` // 数据体
}

func init() {
	host := config.MessageQueue.Host
	port := config.MessageQueue.Port

	Address = net.JoinHostPort(host, port)

	Config = nsq.NewConfig()
	Config.DialTimeout = time.Second * 60
	Config.MsgTimeout = time.Second * 60
	Config.ReadTimeout = time.Second * 60
	Config.WriteTimeout = time.Second * 60
	Config.HeartbeatInterval = time.Second * 10
}

func RunMessageQueueConsumer() ([]*nsq.Consumer, error) {
	consumers := make([]*nsq.Consumer, 0)
	wg := &sync.WaitGroup{}

	wg.Add(100)

	emailConsumer, err := CreateConsumer(TopicSendEmail, ChanelSendEmail, nsq.HandlerFunc(func(message *nsq.Message) error {

		body := SendActivationEmailBody{}

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
			_ = redis.ClientActivationCode.Del(body.Code).Err()
		}

		log.Printf("发送验证码 %s 到 %s\n", body.Code, body.Email)

		return nil
	}))

	if err != nil {
		return consumers, err
	} else {
		consumers = append(consumers, emailConsumer)
	}

	notifyConsumer, err := CreateConsumer(TopicPushNotify, ChanelPushNotify, nsq.HandlerFunc(func(message *nsq.Message) error {

		body := SendNotifyBody{}

		if err := json.Unmarshal(message.Body, &body); err != nil {
			return err
		}

		n := *notify.Notify

		switch body.Event {
		// 发送给所有用户
		case notify.SendNotifyEventSendNotifyToAllUser:
			type SendNotifyPayload struct {
				Title   string `json:"title" valid:"required"`   // 推送的标题
				Content string `json:"content" valid:"required"` // 推送的内容
			}
			b, err := json.Marshal(body.Payload)

			if err != nil {
				return err
			}

			var payload SendNotifyPayload

			if err := json.Unmarshal(b, &payload); err != nil {
				return err
			}

			if err := validator.ValidateStruct(payload); err != nil {
				return err
			}

			if err := n.SendNotifyToAllUser(payload.Title, payload.Content); err != nil {
				return err
			}
			break
		// 发送给指定用户
		case notify.SendNotifyEventSendNotifyToCustomUser:
			type SendNotifyPayload struct {
				UserID  []string `json:"user_id" valid:"required"` // 要指定的推送用户 ID
				Title   string   `json:"title" valid:"required"`   // 推送的标题
				Content string   `json:"content" valid:"required"` // 推送的内容
			}
			b, err := json.Marshal(body.Payload)

			if err != nil {
				return err
			}

			var payload SendNotifyPayload

			if err := json.Unmarshal(b, &payload); err != nil {
				return err
			}

			if err := validator.ValidateStruct(payload); err != nil {
				return err
			}

			if err := n.SendNotifyToCustomUser(payload.UserID, payload.Title, payload.Content); err != nil {
				return err
			}
			break
		// 发送给指定用户，登录异常
		case notify.SendNotifyEventSendNotifyToLoginAbnormalUser:
			b, err := json.Marshal(body.Payload)

			if err != nil {
				return err
			}

			var payload schema.ProfilePublic

			if err := json.Unmarshal(b, &payload); err != nil {
				return err
			}

			if err := validator.ValidateStruct(payload); err != nil {
				return err
			}

			if err := n.SendNotifyToLoginAbnormalUser(payload); err != nil {
				return err
			}
			break
		default:
			return nil
		}

		return nil
	}))

	if err != nil {
		return consumers, err
	} else {
		consumers = append(consumers, notifyConsumer)
	}

	wg.Wait()

	return consumers, nil
}
