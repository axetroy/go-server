package invite

import (
	"fmt"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/jinzhu/gorm"
)

func DeleteByField(field, value string) {
	var (
		err error
		tx  *gorm.DB
	)

	defer func() {
		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}
	}()

	tx = database.Db.Begin()

	raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE %s = '%v'", "invite_history", field, value)

	if err := tx.Exec(raw).Error; err != nil {
		return
	}
}

func DeleteUserById(id string) {
	DeleteByField("id", id)
}
