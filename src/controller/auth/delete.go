package auth

import (
	"github.com/axetroy/go-server/src/service"
)

func DeleteUserByUserName(username string) {
	service.DeleteRowByTable("user", "username", username)
}
