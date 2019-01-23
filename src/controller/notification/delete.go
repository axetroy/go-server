package notification

import (
	"github.com/axetroy/go-server/src/service"
)

func DeleteNotificationById(id string) {
	service.DeleteRowByTable("notification", "id", id)
}

func DeleteNotificationMarkById(id string) {
	service.DeleteRowByTable("notification_mark", "id", id)
}
