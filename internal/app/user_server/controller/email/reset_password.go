// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package email

import (
	"context"
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/captcha"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/email"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/jinzhu/gorm"
	"time"
)

type SendResetPasswordEmailParams struct {
	Email string `json:"email" validate:"required,email,max=255" comment:"邮箱"` // 要发送的邮箱地址
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

		helper.Response(&res, nil, nil, err)
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
	if err = redis.ClientResetCode.Set(context.Background(), code, userInfo.Id, time.Minute*30).Err(); err != nil {
		return
	}

	e, err := email.NewMailer()

	if err != nil {
		return
	}

	// send email
	if err = e.SendForgotPasswordEmail(input.Email, code); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ClientResetCode.Del(context.Background(), code).Err()
		return
	}

	return

}

var SendResetPasswordEmailRouter = router.Handler(func(c router.Context) {
	var (
		input SendResetPasswordEmailParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SendResetPasswordEmail(input)
	})
})
