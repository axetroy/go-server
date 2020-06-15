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
)

// 标记已读
func MarkRead(c helper.Context, notificationID string) (res schema.Response) {
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

		helper.Response(&res, nil, nil, err)
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

var ReadRouter = router.Handler(func(c router.Context) {
	notificationID := c.Param("id")

	c.ResponseFunc(nil, func() schema.Response {
		return MarkRead(helper.NewContext(&c), notificationID)
	})
})

// 批量标记已读
func MarkReadBatch(c helper.Context, notificationIDs []string) (res schema.Response) {
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

		helper.Response(&res, nil, nil, err)
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

	notifications := make([]model.Notification, 0)

	// 先获取通知
	if err = tx.Model(model.Notification{}).Where("id IN (?)", notificationIDs).Find(&notifications).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

loop:
	for _, notificationInfo := range notifications {
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
			continue loop
		}

		if err = tx.Create(&mark).Error; err != nil {
			return
		}
	}

	return
}

type MarkReadBatchParams struct {
	IDs []string `json:"ids"`
}

var MarkReadBatchRouter = router.Handler(func(c router.Context) {
	var params MarkReadBatchParams

	c.ResponseFunc(c.ShouldBindJSON(&params), func() schema.Response {
		return MarkReadBatch(helper.NewContext(&c), params.IDs)
	})
})
