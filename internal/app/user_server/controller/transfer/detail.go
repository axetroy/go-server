// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package transfer

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

func GetDetail(c helper.Context, transferId string) (res schema.Response) {
	var (
		err  error
		tx   *gorm.DB
		data = schema.TransferLog{}
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

		helper.Response(&res, data, nil, err)
	}()

	tx = database.Db.Begin()

	userInfo := model.User{Id: c.Uid}

	if err = tx.Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	log := model.TransferLog{}

	// 联表查询
	// 只能获取自己转给别人的
	sql := GenerateTransferLogSQL(QueryParams{
		Id: &transferId,
	}, 1, false)

	if err = tx.Raw(sql).Scan(&log).Error; err != nil {
		return
	}

	if log.From != c.Uid {
		if log.To != c.Uid {
			// 既不是转账人，也不是收款人, 没有权限获取这条记录
			err = exception.NoPermission
			return
		}
	}

	if err = mapstructure.Decode(log, &data.TransferLogPure); err != nil {
		return
	}

	data.CreatedAt = log.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = log.UpdatedAt.Format(time.RFC3339Nano)
	return
}

var GetDetailRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return GetDetail(helper.NewContext(&c), c.Param("transfer_id"))
	})
})
