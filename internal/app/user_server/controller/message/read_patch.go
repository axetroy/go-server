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

type MarkBatchParams struct {
	IDs []string `json:"ids" valid:"required~请输入 ID 组"`
}

func MarkBatchRead(c helper.Context, input MarkBatchParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
	)

	// 最多 100 个
	if len(input.IDs) > 100 {
		input.IDs = input.IDs[:100]
	}

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

	if err = tx.Model(model.Message{}).Where("id in (?)", input.IDs).Where("uid = ?", c.Uid).Where("read = ?", false).UpdateColumn(model.Message{
		Read:   true,
		ReadAt: &now,
	}).Error; err != nil {
		err = exception.Database
		return
	}

	return
}

var ReadBatchRouter = router.Handler(func(c router.Context) {
	var input MarkBatchParams
	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return MarkBatchRead(helper.NewContext(&c), input)
	})
})
