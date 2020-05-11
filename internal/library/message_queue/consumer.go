// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message_queue

import (
	"github.com/nsqio/go-nsq"
)

// 创建消费者
func CreateConsumer(topic Topic, channel Chanel, handler nsq.Handler) (c *nsq.Consumer, err error) {
	if c, err = nsq.NewConsumer(string(topic), string(channel), Config); err != nil {
		return
	}

	c.AddHandler(handler)

	if err = c.ConnectToNSQD(Address); err != nil {
		return
	}

	return c, nil
}
