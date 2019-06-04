// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"github.com/axetroy/go-server/src/service/database"
)

func DeleteUserByUserName(username string) {
	database.DeleteRowByTable("user", "username", username)
}
