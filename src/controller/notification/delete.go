package notification

import (
	"github.com/axetroy/go-server/src/service"
)

func DeleteNotificationById(id string) {
	service.DeleteRowByTable("notification", "id", id)
}
