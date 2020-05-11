// Copyright 2019 Axetroy. All rights reserved. MIT license.
package email

import (
	"errors"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/captcha"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/email"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type SendResetPasswordEmailParams struct {
	Email string `json:"email" valid:"required~请输入邮箱地址"` // 要发送的邮箱地址
}

func SendResetPasswordEmail(input SendResetPasswordEmailParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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

	userInfo := model.User{
		Email: &input.Email,
	}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 生成重置码
	var code = captcha.GenerateResetCode(userInfo.Id)

	// set activationCode to redis
	if err = redis.ClientResetCode.Set(code, userInfo.Id, time.Minute*30).Err(); err != nil {
		return
	}

	e := email.NewMailer()

	// send email
	if err = e.SendForgotPasswordEmail(input.Email, code); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ClientResetCode.Del(code).Err()
		return
	}

	return

}

func SendResetPasswordEmailRouter(c *gin.Context) {
	var (
		input SendResetPasswordEmailParams
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

	res = SendResetPasswordEmail(input)
}
