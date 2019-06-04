// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/service/redis"
	"github.com/axetroy/go-server/src/util"
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
		err = exception.InvalidParams
		return
	}

	if uid, err = redis.ResetCodeClient.Get(input.Code).Result(); err != nil {
		err = exception.InvalidResetCode
		return
	}

	tx = database.Db.Begin()

	userInfo := model.User{Id: uid}

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
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

func ResetPasswordRouter(context *gin.Context) {
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
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = ResetPassword(input)
}
