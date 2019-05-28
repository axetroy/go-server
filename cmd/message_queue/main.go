package main

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/src/message_queue"
	"github.com/axetroy/go-server/src/message_queue/consumer"
	"github.com/axetroy/go-server/src/service"
	"github.com/axetroy/go-server/src/service/email"
	"github.com/nsqio/go-nsq"
	"log"
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)

	_, err := consumer.CreateConsumer(message_queue.TopicSendEmail, message_queue.ChanelSendEmail, nsq.HandlerFunc(func(message *nsq.Message) error {

		body := message_queue.SendEmailBody{}

		if err := json.Unmarshal(message.Body, &body); err != nil {
			return err
		}

		mailer := email.NewMailer()

		// 发送邮件
		if err := mailer.SendActivationEmail(body.Email, body.Code); err != nil {
			// 邮件没发出去的话，删除 redis 的 key
			_ = service.RedisActivationCodeClient.Del(body.Code).Err()
		}

		fmt.Printf("发送验证码 %s 到 %s\n", body.Code, body.Email)

		return nil
	}))

	if err != nil {
		log.Panic(err)
	}

	wg.Wait()
}
