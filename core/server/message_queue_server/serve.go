// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message_queue_server

import (
	"context"
	"github.com/axetroy/go-server/core/config"
	"github.com/axetroy/go-server/core/message_queue"
	"github.com/axetroy/go-server/core/service/database"
	"github.com/nsqio/go-nsq"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve() error {
	var (
		c *nsq.Consumer
	)

	go func() {
		if ctx, err := message_queue.RunMessageQueueConsumer(); err != nil {
			log.Fatal(err)
		} else {
			c = ctx
		}
	}()

	log.Println("Listening message queue")

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	config.Common.Exiting = true

	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if c != nil {
		c.Stop()

		_ = c.DisconnectFromNSQD(message_queue.Address)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("Timeout of 5 seconds.")
	}

	_ = database.Db.Close()

	log.Println("Server exiting")

	return nil
}
