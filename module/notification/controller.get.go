// Copyright 2019 Axetroy. All rights reserved. MIT license.
package notification

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/notification/notification_model"
	"github.com/axetroy/go-server/module/notification/notification_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

// Get notification detail
func Get(context schema.Context, id string) (res schema.Response) {
	var (
		err  error
		data notification_schema.Notification
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	notificationInfo := notification_model.Notification{}
	NotificationMark := notification_model.NotificationMark{Id: id}

	sql := fmt.Sprintf("LEFT JOIN notification_mark ON notification_mark.id = notification.id AND notification.id = '%s'", notificationInfo.Id)

	if err = tx.Table(notificationInfo.TableName()).Joins(sql).Last(&notificationInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrNoData
		}
		return
	}

	if er := mapstructure.Decode(notificationInfo, &data.NotificationPure); er != nil {
		err = er
		return
	}

	if err = tx.Where(&NotificationMark).Last(&NotificationMark).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			data.NotificationPure.Read = false
			err = nil
		} else {
			return
		}
	} else {
		data.NotificationPure.Read = NotificationMark.Read
		data.ReadAt = NotificationMark.CreatedAt.Format(time.RFC3339Nano)
	}

	data.CreatedAt = notificationInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = notificationInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

// GetRouter get notification detail router
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

	id := ctx.Param("id")

	res = Get(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}
