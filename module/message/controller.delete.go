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

func DeleteMessageById(id string) {
	database.DeleteRowByTable("message", "id", id)
}

func DeleteByAdmin(context schema.Context, messageId string) (res schema.Response) {
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
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	adminInfo := admin_model.Admin{Id: context.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	messageInfo := message_model.Message{
		Id: messageId,
	}

	if err = tx.First(&messageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrMessageNotExist
			return
		}
		return
	}

	if err = tx.Delete(message_model.Message{Id: messageInfo.Id}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(messageInfo, &data.MessagePure); err != nil {
		return
	}

	data.CreatedAt = messageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = messageInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func DeleteByAdminRouter(ctx *gin.Context) {
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

	res = DeleteByAdmin(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}

func DeleteByUser(context schema.Context, messageId string) (res schema.Response) {
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
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	messageInfo := message_model.Message{
		Id:  messageId,
		Uid: context.Uid,
	}

	if err = tx.First(&messageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrMessageNotExist
			return
		}
		return
	}

	if err = tx.Delete(message_model.Message{Id: messageInfo.Id}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(messageInfo, &data.MessagePure); err != nil {
		return
	}

	data.CreatedAt = messageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = messageInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func DeleteByUserRouter(ctx *gin.Context) {
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

	res = DeleteByUser(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}
