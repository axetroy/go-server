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

// Query params
type Query struct {
	schema.Query
}

// GetList get notification list
func GetNotificationListByUser(c helper.Context, input Query) (res schema.Response) {
	var (
		err  error
		data = make([]schema.Notification, 0)
		meta = &schema.Meta{}
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

		helper.Response(&res, data, meta, err)
	}()

	query := input.Query

	query.Normalize()

	tx = database.Db.Begin()

	list := make([]model.Notification, 0)

	filter := map[string]interface{}{}

	if err = query.Order(tx.Limit(query.Limit).Offset(query.Limit * query.Page).Where(filter)).Find(&list).Error; err != nil {
		return
	}

	var total int64

	if err = tx.Model(&model.Notification{}).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.Notification{}
		if er := mapstructure.Decode(v, &d.NotificationPure); er != nil {
			err = er
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)

		// 查询用户是否已读通知
		mark := model.NotificationMark{
			Id:  v.Id,
			Uid: c.Uid,
		}

		if err = tx.Last(&mark).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				d.Read = false
				d.ReadAt = ""
				err = nil
			} else {
				break
			}
		} else {
			d.Read = true
			d.ReadAt = mark.CreatedAt.Format(time.RFC3339Nano)
		}

		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}

// GetListRouter get list router
var GetNotificationListByUserRouter = router.Handler(func(c router.Context) {
	var input Query

	c.ResponseFunc(c.ShouldBindQuery(&input), func() schema.Response {
		return GetNotificationListByUser(helper.NewContext(&c), input)
	})
})
