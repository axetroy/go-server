// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message_queue_server

import (
	"github.com/axetroy/go-server/src/message_queue"
)

func Serve() error {
	return message_queue.RunMessageQueueConsumer()
}
