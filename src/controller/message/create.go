// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateMessageParams struct {
	Uid     string `json:"uid" valid:"required~请添加用户ID"`
	Title   string `json:"title" valid:"required~请填写消息标题"`
	Content string `json:"content" valid:"required~请填写消息内容"`
}

func Create(context controller.Context, input CreateMessageParams) (res schema.Response) {
	var (
		err          error
		data         schema.Message
		tx           *gorm.DB
		isValidInput bool
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		err = exception.WrapValidatorError(err)
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Id: context.Uid,
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

func CreateRouter(c *gin.Context) {
	var (
		input CreateMessageParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Create(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, input)
}
