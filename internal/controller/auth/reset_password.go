// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"errors"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/util"
	"github.com/axetroy/go-server/internal/validator"
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
		err error
		tx  *gorm.DB
		uid string // 重置码对应的uid
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

		helper.Response(&res, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if uid, err = redis.ClientResetCode.Get(input.Code).Result(); err != nil {
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
	if err = redis.ClientResetCode.Del(input.Code).Err(); err != nil {
		return
	}

	return
}

func ResetPasswordRouter(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = ResetPassword(input)
}
