// Copyright 2019 Axetroy. All rights reserved. MIT license.
package invite

import (
	"github.com/axetroy/go-server/src/service/database"
)

func DeleteById(id string) {
	database.DeleteRowByTable("invite_history", "id", id)
}
