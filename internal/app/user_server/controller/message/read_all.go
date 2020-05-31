// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"time"
)

func MarkAllRead(c helper.Context) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, true, nil, err)
	}()

	tx = database.Db.Begin()

	now := time.Now()

	if err = tx.Model(model.Message{}).Where("uid = ?", c.Uid).Where("read = ?", false).UpdateColumn(model.Message{
		Read:   true,
		ReadAt: &now,
	}).Error; err != nil {
		err = exception.Database
		return
	}

	return
}

var ReadAllRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return MarkAllRead(helper.NewContext(&c))
	})
})
