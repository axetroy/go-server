package consumer

import (
	"github.com/axetroy/go-server/src/message_queue"
	"github.com/nsqio/go-nsq"
)

// 创建消费者
func CreateConsumer(topic message_queue.Topic, channel message_queue.Chanel, handler nsq.Handler) (c *nsq.Consumer, err error) {
	c, err = nsq.NewConsumer(string(topic), string(channel), message_queue.Config)

	if err != nil {
		return
	}

	c.AddHandler(handler)

	err = c.ConnectToNSQD(message_queue.Address)

	if err != nil {
		return
	}

	return c, nil
}
