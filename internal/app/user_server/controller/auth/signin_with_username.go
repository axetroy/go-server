// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package auth

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/authentication"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

type SignInParams struct {
	Account  string `json:"account" validate:"required,max=36" comment:"帐号"`
	Password string `json:"password" validate:"required,max=32" comment:"密码"`
	Duration *int64 `json:"duration" validate:"omitempty,number,gt=0" comment:"有效时间"`
}

// 普通帐号登陆
func SignIn(c helper.Context, input SignInParams) (res schema.Response) {
	var (
		err  error
		data = &schema.ProfileWithToken{}
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

		helper.Response(&res, data, nil, err)
	}()

	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	userInfo := model.User{
		Password: util.GeneratePassword(input.Password),
	}

	if validator.IsPhone(input.Account) {
		// 用手机号登陆
		userInfo.Phone = &input.Account
	} else if validator.IsEmail(input.Account) {
		// 用邮箱登陆
		userInfo.Email = &input.Account
	} else {
		// 用用户名
		userInfo.Username = input.Account
	}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Preload("Wechat").Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
		return
	}

	// 检查用户登录状态
	go func() {
		if er := message_queue.PublishCheckUserLogin(userInfo.Id); er != nil {
			log.Println("检查用户状态失败", c.Uid)
		}
	}()

	if err = userInfo.CheckStatusValid(); err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	if userInfo.WechatOpenID != nil {
		if err = mapstructure.Decode(userInfo.Wechat, &data.Wechat); err != nil {
			return
		}
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	var duration time.Duration

	if input.Duration != nil {
		duration = time.Duration(*input.Duration * int64(time.Second))
	} else {
		duration = time.Hour * 6
	}

	// generate token
	if t, er := authentication.Gateway(false).Generate(userInfo.Id, duration); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	loginLog := model.LoginLog{
		Uid:     userInfo.Id,                       // 用户ID
		Type:    model.LoginLogTypeUserName,        // 默认用户名登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  c.UserAgent,                       // 用户的 userAgent
		LastIp:  c.Ip,                              // 用户的IP
	}

	if err = tx.Create(&loginLog).Error; err != nil {
		return
	}

	return
}

var SignInRouter = router.Handler(func(c router.Context) {
	var (
		input SignInParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignIn(helper.NewContext(&c), input)
	})
})
