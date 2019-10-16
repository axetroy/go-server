// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message_queue

import (
	"errors"
	"github.com/nsqio/go-nsq"
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
