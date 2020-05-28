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
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

type CreateParams struct {
	Title   string  `json:"title" valid:"required~请输入公告标题"`   // 公告标题
	Content string  `json:"content" valid:"required~请输入公告内容"` // 公告内容
	Note    *string `json:"note"`                             // 备注
}

func Create(c helper.Context, input CreateParams) (res schema.Response) {
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

		go func() {
			if err != nil {
				if er := message_queue.PublishSystemNotify(data.Id); er != nil {
					log.Println("加入推送队列失败:", err.Error())
				} else {
					log.Println("假如队列成功...")
				}
			}
		}()

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
		Author:  adminInfo.Id,
		Title:   input.Title,
		Content: input.Content,
		Note:    input.Note,
	}

	if err = tx.Create(&notificationInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(notificationInfo, &data.NotificationPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = notificationInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = notificationInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

var CreateRouter = router.Handler(func(c router.Context) {
	var (
		input CreateParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Create(helper.NewContext(&c), input)
	})
})
