// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/notification/notification_model"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// MarkRead mark notification as read
func MarkRead(context schema.Context, notificationID string) (res schema.Response) {
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
				err = common_error.ErrUnknown
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
			res.Data = false
			res.Message = err.Error()
		} else {
			res.Data = true
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	userInfo := user_model.User{
		Id: context.Uid,
	}

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		// 没有找到用户
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
		}
		return
	}

	notificationInfo := notification_model.Notification{
		Id: notificationID,
	}

	// 先获取通知
	if err = tx.Where(&notificationInfo).Last(&notificationInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrNoData
		}
		return
	}

	mark := notification_model.NotificationMark{
		Id:  notificationInfo.Id,
		Uid: context.Uid,
	}

	// 再确认以读表有没有这个用户的已读记录
	if err = tx.Where(&mark).Last(&mark).Error; err != nil {
		// 如果没找到这条记录，则说明没有创建
		// 继续下面的页面
		if err == gorm.ErrRecordNotFound {
			err = nil
		} else {
			return
		}
	} else {
		// 通知已读
		return
	}

	if err = tx.Create(&mark).Error; err != nil {
		return
	}

	return
}

// ReadRouter read this notification router
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

	notificationID := ctx.Param("id")

	res = MarkRead(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, notificationID)
}
