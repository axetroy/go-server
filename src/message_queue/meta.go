package message_queue

import (
	"github.com/axetroy/go-server/src/config"
	"github.com/nsqio/go-nsq"
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
