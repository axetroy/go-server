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
	"github.com/axetroy/go-server/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type ResetPasswordParams struct {
	Code        string `json:"code" valid:"required~请输入激活码"`
	NewPassword string `json:"new_password" valid:"required~请输入新密码"`
}

func ResetPassword(input ResetPasswordParams) (res schema.Response) {
	var (
		err          error
		tx           *gorm.DB
		uid          string // 重置码对应的uid
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
			res.Data = false
		} else {
			res.Status = schema.StatusSuccess
			res.Data = true
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.ErrInvalidParams
		return
	}

	if uid, err = redis.ResetCodeClient.Get(input.Code).Result(); err != nil {
		err = user_error.ErrInvalidResetCode
		return
	}

	tx = database.Db.Begin()

	userInfo := user_model.User{Id: uid}

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	// 更新密码
	tx.Model(&userInfo).Update("password", util.GeneratePassword(input.NewPassword))

	// delete reset code from redis
	if err = redis.ResetCodeClient.Del(input.Code).Err(); err != nil {
		return
	}

	// TODO: 安全起见，发送一封邮件/短信告知用户
	return
}

func ResetPasswordRouter(ctx *gin.Context) {
	var (
		input ResetPasswordParams
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

	res = ResetPassword(input)
}
