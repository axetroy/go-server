// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/module/log/log_model"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type SignInParams struct {
	Account  string  `json:"account" valid:"required~请输入登陆账号"`
	Password string  `json:"password" valid:"required~请输入密码"`
	Code     *string `json:"code"` // 手机验证码
}

func SignIn(context schema.Context, input SignInParams) (res schema.Response) {
	var (
		err          error
		data         = &user_schema.ProfileWithToken{}
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
				err = common_error.ErrUnknown
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}

	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = common_error.ErrInvalidParams
		return
	}

	userInfo := user_model.User{Password: util.GeneratePassword(input.Password)}

	if govalidator.Matches(input.Account, "^/d+$") && input.Code != nil { // 如果是手机号, 并且传入了code字段
		// TODO: 这里手机登陆应该用验证码作为校验
		userInfo.Phone = &input.Account
	} else if govalidator.IsEmail(input.Account) { // 如果是邮箱的话
		userInfo.Email = &input.Account
	} else {
		userInfo.Username = input.Account // 其他则为用户名
	}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrInvalidAccountOrPassword
		}
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	if t, er := token.Generate(userInfo.Id, false); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	log := log_model.LoginLog{
		Uid:     userInfo.Id,
		Type:    log_model.LoginLogTypeUserName,        // 默认用户名登陆
		Command: log_model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  context.UserAgent,
		LastIp:  context.Ip,
	}

	if err = tx.Create(&log).Error; err != nil {
		return
	}

	return
}

func SignInRouter(ctx *gin.Context) {
	var (
		input SignInParams
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
		err = common_error.ErrInvalidParams
		return
	}

	res = SignIn(schema.Context{
		UserAgent: ctx.GetHeader("user-agent"),
		Ip:        ctx.ClientIP(),
	}, input)
}
