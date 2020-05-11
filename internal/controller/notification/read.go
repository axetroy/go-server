// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification

import (
	"errors"
	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/helper"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// MarkRead mark notification as read
func MarkRead(c controller.Context, notificationID string) (res schema.Response) {
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

		helper.Response(&res, nil, err)
	}()

	tx = database.Db.Begin()

	userInfo := model.User{
		Id: c.Uid,
	}

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		// 没有找到用户
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	notificationInfo := model.Notification{
		Id: notificationID,
	}

	// 先获取通知
	if err = tx.Where(&notificationInfo).Last(&notificationInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	mark := model.NotificationMark{
		Id:  notificationInfo.Id,
		Uid: c.Uid,
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
func ReadRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	notificationID := c.Param("id")

	res = MarkRead(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, notificationID)
}
