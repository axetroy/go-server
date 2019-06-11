// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/message/message_model"
	"github.com/axetroy/go-server/module/message/message_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func MarkRead(context schema.Context, id string) (res schema.Response) {
	var (
		err  error
		data message_schema.Message
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
				err = exception.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	MessageInfo := message_model.Message{
		Id:  id,
		Uid: context.Uid,
	}

	if err = tx.Where(&MessageInfo).Last(&MessageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.ErrNoData
		}
		return
	}

	if er := mapstructure.Decode(MessageInfo, &data.MessagePure); er != nil {
		err = er
		return
	}

	data.CreatedAt = MessageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = MessageInfo.UpdatedAt.Format(time.RFC3339Nano)

	now := time.Now()

	if err = tx.Model(&MessageInfo).UpdateColumn(message_model.Message{
		Read:   true,
		ReadAt: &now,
	}).Error; err != nil {
		return
	}

	return
}

func ReadRouter(ctx *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	id := ctx.Param(ParamsIdName)

	res = MarkRead(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}
