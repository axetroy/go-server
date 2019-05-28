package admin

import (
	"github.com/axetroy/go-server/src/service/database"
)

func DeleteAdminByAccount(account string) {
	database.DeleteRowByTable("admin", "username", account)
}
