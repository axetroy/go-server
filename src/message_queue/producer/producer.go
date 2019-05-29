package producer

import (
	"fmt"
	"github.com/axetroy/go-server/src/message_queue"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
)

var (
	producer *nsq.Producer
)

func init() {
	host := os.Getenv("MSG_QUEUE_SERVER")

	if host == "" {
		host = "127.0.0.1"
	}

	port := os.Getenv("MSG_QUEUE_PORT")

	if port == "" {
		port = "4150"
	}

	addr := host + ":" + port

	// TODO: 断线重连机制
	err := CreateProducer(addr)

	if err != nil {
		log.Panic(err)
	}
}

// 初始化生产者
func CreateProducer(address string) (err error) {
	producer, err = nsq.NewProducer(address, nsq.NewConfig())
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
