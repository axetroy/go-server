package address

import (
	"github.com/axetroy/go-server/src/service"
)

func DeleteAddressById(id string) {
	service.DeleteRowByTable("address", "id", id)
}
