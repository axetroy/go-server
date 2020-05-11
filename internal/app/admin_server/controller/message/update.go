// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type UpdateParams struct {
	Title   *string `json:"title"`   // 消息标题
	Content *string `json:"content"` // 消息内容
}

func Update(c helper.Context, messageId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.Message
		tx           *gorm.DB
		shouldUpdate bool
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
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, err)
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

	messageInfo := model.Message{
		Id: messageId,
	}

	if err = tx.First(&messageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.MessageNotExist
			return
		}
		return
	}

	updateModel := model.Message{}

	if input.Title != nil {
		shouldUpdate = true
		updateModel.Title = *input.Title
	}

	if input.Content != nil {
		shouldUpdate = true
		updateModel.Content = *input.Content
	}

	if shouldUpdate {
		if err = tx.Model(&messageInfo).Updates(&updateModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.MessageNotExist
				return
			}
			return
		}
	}

	if err = mapstructure.Decode(messageInfo, &data.MessagePure); err != nil {
		return
	}

	data.CreatedAt = messageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = messageInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	id := c.Param(ParamsIdName)

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Update(helper.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, id, input)
}
