// Copyright 2019 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

var (
	PayPasswordHeader = "X-Pay-Password"
)

// 交易密码的验证中间件
func AuthPayPassword(ctx *gin.Context) {
	var (
		err error
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

		// 如果有报错的话，那么不会进入到路由中
		if err != nil {
			ctx.JSON(http.StatusOK, schema.Response{
				Status:  schema.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})

			// 中断后面的路由器执行
			ctx.Abort()

			return
		}
	}()

	payPassword := ctx.GetHeader(PayPasswordHeader)

	if len(payPassword) == 0 {
		err = common_error.ErrRequirePayPassword
		return
	}

	uid := ctx.GetString(ContextUidField)

	if uid == "" {
		err = common_error.ErrUserNotLogin
		return
	}

	userInfo := user_model.User{Id: uid}

	if err = database.Db.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
		}
		return
	}

	if userInfo.PayPassword == nil {
		err = common_error.ErrPayPasswordNotSet
		return
	}

	// 校验密码是否正确
	if *userInfo.PayPassword != util.GeneratePassword(payPassword) {
		err = common_error.ErrInvalidPassword
		return
	}

}
