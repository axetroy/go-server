package producer

import (
	"fmt"
	"github.com/axetroy/go-server/src/message_queue"
	"github.com/nsqio/go-nsq"
	"log"
)

var (
	producer *nsq.Producer
)

func init() {
	// TODO: 断线重连机制
	err := CreateProducer(message_queue.Address)

	if err != nil {
		log.Panic(err)
	}
}

// 初始化生产者
func CreateProducer(address string) (err error) {
	producer, err = nsq.NewProducer(address, message_queue.Config)
	fmt.Printf("连接队列: %s\n", address)
	return
}

// 发布消息
func Publish(topic message_queue.Topic, message []byte) error {
	var err error

	if producer != nil {
		if len(message) == 0 { //不能发布空串，否则会导致error
			return nil
		}

		err = producer.Publish(string(topic), message) // 发布消息

		// 如果发送失败，则重新连接
		if err != nil {
			err = CreateProducer("")
		}

		return err
	}

	return fmt.Errorf("producer is nil %v", err)
}
