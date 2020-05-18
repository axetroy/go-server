// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package notification

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

func DeleteNotificationById(id string) {
	database.DeleteRowByTable("notification", "id", id)
}

func DeleteNotificationMarkById(id string) {
	database.DeleteRowByTable("notification_mark", "id", id)
}

func Delete(c helper.Context, notificationId string) (res schema.Response) {
	var (
		err  error
		data schema.Notification
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

	adminInfo := model.Admin{Id: c.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	notificationInfo := model.Notification{
		Id: notificationId,
	}

	if err = tx.First(&notificationInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NotificationNotExist
			return
		}
		return
	}

	if err = tx.Delete(model.Notification{Id: notificationInfo.Id}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(notificationInfo, &data.NotificationPure); err != nil {
		return
	}

	data.CreatedAt = notificationInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = notificationInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var DeleteRouter = router.Handler(func(c router.Context) {
	id := c.Param("id")

	c.ResponseFunc(nil, func() schema.Response {
		return Delete(helper.NewContext(&c), id)
	})
})
