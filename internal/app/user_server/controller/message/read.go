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
	"github.com/mitchellh/mapstructure"
	"time"
)

func MarkRead(c helper.Context, id string) (res schema.Response) {
	var (
		err  error
		data schema.Message
		tx   *gorm.DB
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

	MessageInfo := model.Message{
		Id:  id,
		Uid: c.Uid,
	}

	if err = tx.Where(&MessageInfo).Last(&MessageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if er := mapstructure.Decode(MessageInfo, &data.MessagePure); er != nil {
		err = er
		return
	}

	now := time.Now()

	if err = tx.Model(&MessageInfo).UpdateColumn(model.Message{
		Read:   true,
		ReadAt: &now,
	}).Error; err != nil {
		return
	}

	nowStr := now.Format(time.RFC3339Nano)

	data.Read = true
	data.ReadAt = &nowStr
	data.CreatedAt = MessageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = MessageInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var ReadRouter = router.Handler(func(c router.Context) {
	id := c.Param("message_id")

	c.ResponseFunc(nil, func() schema.Response {
		return MarkRead(helper.NewContext(&c), id)
	})
})
