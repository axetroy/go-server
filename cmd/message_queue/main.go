// Copyright 2019 Axetroy. All rights reserved. MIT license.
package main

import (
	"github.com/axetroy/go-server/module/message_queue"
)

func main() {
	message_queue.RunMessageQueueConsumer()
}
