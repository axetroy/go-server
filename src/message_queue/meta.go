package message_queue

import (
	"github.com/nsqio/go-nsq"
	"os"
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
	host := os.Getenv("MSG_QUEUE_SERVER")

	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("MSG_QUEUE_PORT")

	if port == "" {
		port = "4150"
	}

	Address = host + ":" + port

	Config = nsq.NewConfig()
	Config.DialTimeout = time.Second * 5
	Config.MsgTimeout = time.Second * 10
	Config.ReadTimeout = time.Second * 10
	Config.WriteTimeout = time.Second * 10
}
