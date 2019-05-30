package main

import (
	"github.com/axetroy/go-server/src/message_queue"
)

func main() {
	message_queue.RunMessageQueueConsumer()
}
