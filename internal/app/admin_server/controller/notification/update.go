// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package notification

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type UpdateParams struct {
	Title   *string `json:"title" validate:"omitempty,max=32" comment:"标题"` // 公告标题
	Content *string `json:"content" validate:"omitempty" comment:"内容"`      // 公告内容
	Note    *string `json:"note" validate:"omitempty,max=255" comment:"备注"` // 备注
}

func Update(c helper.Context, notificationId string, input UpdateParams) (res schema.Response) {
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

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.Where(&adminInfo).First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	notificationInfo := model.Notification{
		Id: notificationId,
	}

	if err = tx.Where(&notificationInfo).Last(&notificationInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	updateModel := model.Notification{}

	if input.Title != nil && len(*input.Title) != 0 {
		updateModel.Title = *input.Title
	}

	if input.Content != nil && len(*input.Content) != 0 {
		updateModel.Content = *input.Content
	}

	if input.Note != nil {
		updateModel.Note = input.Note
	}

	if err = tx.Model(&notificationInfo).Updates(&updateModel).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(notificationInfo, &data.NotificationPure); err != nil {
		return
	}

	data.CreatedAt = notificationInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = notificationInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var UpdateRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Update(helper.NewContext(&c), c.Param("id"), input)
	})
})
