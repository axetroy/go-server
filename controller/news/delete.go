package news

import (
	"fmt"
	"github.com/axetroy/go-server/orm"
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
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}
	}()

	tx = orm.DB.Begin()

	raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE %s = '%v'", "news", field, value)

	if err = tx.Exec(raw).Error; err != nil {
		return
	}
}

func DeleteNewsById(id string) {
	DeleteByField("id", id)
}
