package notification

import (
	"fmt"
	"github.com/axetroy/go-server/orm"
	"github.com/go-xorm/xorm"
)

func DeleteByField(field, value string) {
	var (
		err     error
		session *xorm.Session
		tx      bool
	)

	defer func() {
		if tx {
			if err != nil {
				_ = session.Rollback()
			} else {
				_ = session.Commit()
			}
		}
	}()

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE %s = '%v'", "notification", field, value)

	if _, err := session.Exec(raw); err != nil {
		return
	} else {

	}
}

func DeleteNotificationById(id string) {
	DeleteByField("id", id)
}
