// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message

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
	"time"
)

type CreateMessageParams struct {
	Uid     string `json:"uid" validate:"required,max=32" comment:"用户 ID "`
	Title   string `json:"title" validate:"required,max=32" comment:"标题"`
	Content string `json:"content" validate:"required" comment:"内容"`
}

func Create(c helper.Context, input CreateMessageParams) (res schema.Response) {
	var (
		err  error
		data schema.Message
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

		// 推送用户
		if err != nil {
			// 把它加入到队列中
			_ = message_queue.PublishUserMessage(data.Id)
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

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	if !adminInfo.IsSuper {
		err = exception.AdminNotSuper
		return
	}

	userInfo := model.User{
		Id: input.Uid,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		// 没有找到用户
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	MessageInfo := model.Message{
		Uid:     input.Uid,
		Title:   input.Title,
		Content: input.Content,
		Status:  model.MessageStatusActive,
	}

	if err = tx.Create(&MessageInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(MessageInfo, &data.MessagePure); er != nil {
		err = er
		return
	}

	data.CreatedAt = MessageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = MessageInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

var CreateRouter = router.Handler(func(c router.Context) {
	var (
		input CreateMessageParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Create(helper.NewContext(&c), input)
	})
})
