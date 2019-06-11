// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/service/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type ActivationParams struct {
	Code string `json:"code" valid:"required~请输入激活码;"`
}

func Activation(input ActivationParams) (res schema.Response) {
	var (
		err          error
		tx           *gorm.DB
		uid          string // 激活码对应的用户ID
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
				err = exception.ErrUnknown
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
			res.Message = err.Error()
		} else {
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.ErrInvalidParams
		return
	}

	if uid, err = redis.ActivationCodeClient.Get(input.Code).Result(); err != nil {
		err = exception.ErrInvalidActiveCode
		return
	}

	tx = database.Db.Begin()

	userInfo := user_model.User{Id: uid}

	if err = tx.Where(&userInfo).Find(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	// 如果用户的状态不是未激活的话
	if userInfo.Status != user_model.UserStatusInactivated {
		err = exception.ErrUserHaveActive
		return
	}

	// 更新激活状态
	tx.Model(&userInfo).Update("status", user_model.UserStatusInit)

	// delete code from redis
	if err = redis.ActivationCodeClient.Del(input.Code).Err(); err != nil {
		return
	}
	return
}

func ActivationRouter(ctx *gin.Context) {
	var (
		input ActivationParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = Activation(input)
}
