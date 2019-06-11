// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_model"
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

// Get Message detail
func Get(context schema.Context, id string) (res schema.Response) {
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

	if err = tx.Last(&MessageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.ErrNoData
		}
		return
	}

	if err = mapstructure.Decode(MessageInfo, &data.MessagePure); err != nil {
		return
	}

	if MessageInfo.ReadAt != nil {
		readAt := MessageInfo.ReadAt.Format(time.RFC3339Nano)
		data.Read = true
		data.ReadAt = &readAt
	}

	data.CreatedAt = MessageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = MessageInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

// Get Message detail
func GetByAdmin(context schema.Context, id string) (res schema.Response) {
	var (
		err  error
		data message_schema.MessageAdmin
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

	adminInfo := admin_model.Admin{
		Id: context.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	MessageInfo := message_model.Message{
		Id: id,
	}

	if err = tx.Last(&MessageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.ErrNoData
		}
		return
	}

	if err = mapstructure.Decode(MessageInfo, &data.MessagePureAdmin); err != nil {
		return
	}

	if MessageInfo.ReadAt != nil {
		readAt := MessageInfo.ReadAt.Format(time.RFC3339Nano)
		data.Read = true
		data.ReadAt = &readAt
	}

	data.CreatedAt = MessageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = MessageInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

// GetRouter get Message detail router
func GetRouter(ctx *gin.Context) {
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

	res = Get(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}

// 管理员获取个人消息详情
func GetAdminRouter(ctx *gin.Context) {
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

	res = GetByAdmin(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}
