// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris/v12"
	"net/http"
)

var (
	PayPasswordHeader = "X-Pay-Password"
	SignatureHeader   = "X-Signature"
)

// 交易密码的验证中间件
func AuthPayPassword(c *gin.Context) {
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
				err = exception.Unknown
			}
		}

		// 如果有报错的话，那么不会进入到路由中
		if err != nil {
			c.JSON(http.StatusOK, schema.Response{
				Status:  schema.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})

			// 中断后面的路由器执行
			c.Abort()

			return
		}
	}()

	payPassword := c.GetHeader(PayPasswordHeader)

	if len(payPassword) == 0 {
		err = exception.RequirePayPassword
		return
	}

	uid := c.GetString(ContextUidField)

	if uid == "" {
		err = exception.UserNotLogin
		return
	}

	userInfo := model.User{Id: uid}

	if err = database.Db.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if userInfo.PayPassword == nil {
		err = exception.RequirePayPasswordSet
		return
	}

	// 校验密码是否正确
	if *userInfo.PayPassword != util.GeneratePassword(payPassword) {
		err = exception.InvalidPassword
		return
	}

}

// 交易密码的验证中间件
func AuthPayPasswordNew(c iris.Context) {
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
				err = exception.Unknown
			}
		}

		// 如果有报错的话，那么不会进入到路由中
		if err != nil {
			_, _ = c.JSON(schema.Response{
				Status:  schema.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})
			return
		} else {
			c.Next()
		}
	}()

	payPassword := c.GetHeader(PayPasswordHeader)

	if len(payPassword) == 0 {
		err = exception.RequirePayPassword
		return
	}

	uid := c.Values().GetString(ContextUidField)

	if uid == "" {
		err = exception.UserNotLogin
		return
	}

	userInfo := model.User{Id: uid}

	if err = database.Db.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if userInfo.PayPassword == nil {
		err = exception.RequirePayPasswordSet
		return
	}

	// 校验密码是否正确
	if *userInfo.PayPassword != util.GeneratePassword(payPassword) {
		err = exception.InvalidPassword
		return
	}
}
