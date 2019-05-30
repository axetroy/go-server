package producer

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/src/message_queue"
	"github.com/nsqio/go-nsq"
	"log"
)

var (
	producer *nsq.Producer
)

func init() {
	err := CreateProducer(message_queue.Address)

	if err != nil {
		log.Panic(err)
	}
}

// 初始化生产者
func CreateProducer(address string) (err error) {
	producer, err = nsq.NewProducer(address, message_queue.Config)
	return
}

// 发布消息
func Publish(topic message_queue.Topic, message []byte) error {
	var err error

	if producer != nil {
		if len(message) == 0 { //不能发布空串，否则会导致error
			return errors.New("message can not be empty")
		}

		err = producer.Publish(string(topic), message) // 发布消息

		return err
	} else {
		err = errors.New("未连接队列")
	}

	return fmt.Errorf("producer is nil %v", err)
}
