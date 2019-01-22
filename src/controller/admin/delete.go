package admin

import (
	"github.com/axetroy/go-server/src/service"
)

func DeleteAdminByAccount(account string) {
	service.DeleteRowByTable("admin", "username", account)
}
