// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"github.com/axetroy/go-server/internal/service/database"
)

func DeleteUserByUserName(username string) {
	database.DeleteRowByTable("user", "username", username)
}

func DeleteUserByUid(uid string) {
	database.DeleteRowByTable("user", "id", uid)
	database.DeleteRowByTable("wechat_open_id", "uid", uid)
}
