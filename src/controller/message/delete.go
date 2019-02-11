package message

import (
	"github.com/axetroy/go-server/src/service"
)

func DeleteMessageById(id string) {
	service.DeleteRowByTable("message", "id", id)
}
