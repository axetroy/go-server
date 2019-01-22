package news

import (
	"github.com/axetroy/go-server/src/service"
)

func DeleteNewsById(id string) {
	service.DeleteRowByTable("news", "id", id)
}
