package consumer

import (
	"github.com/axetroy/go-server/src/message_queue"
	"github.com/nsqio/go-nsq"
)

// 创建消费者
func CreateConsumer(topic message_queue.Topic, channel message_queue.Chanel, handler nsq.Handler) (c *nsq.Consumer, err error) {
	config := nsq.NewConfig()

	c, err = nsq.NewConsumer(string(topic), string(channel), config)

	if err != nil {
		return
	}

	c.AddHandler(handler)

	// TODO: 从环境变量中读取
	// TODO: 断线重连机制
	err = c.ConnectToNSQD("127.0.0.1:4150")

	if err != nil {
		return
	}

	return c, nil
}
