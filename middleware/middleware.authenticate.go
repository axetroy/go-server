// Copyright 2019 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ContextUidField = "uid"
)

// Token 验证中间件
func Authenticate(isAdmin bool) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var (
			err         error
			tokenString string
		)
		defer func() {
			if err != nil {
				ctx.JSON(http.StatusOK, schema.Response{
					Message: err.Error(),
					Data:    nil,
				})
				ctx.Abort()
			}
		}()

		if s, isExist := ctx.GetQuery(token.AuthField); isExist == true {
			tokenString = s
			return
		} else {
			tokenString = ctx.GetHeader(token.AuthField)

			if len(tokenString) == 0 {
				if s, er := ctx.Cookie(token.AuthField); er != nil {
					err = token.ErrInvalidToken
					return
				} else {
					tokenString = s
				}
			}
		}

		if claims, er := token.Parse(tokenString, isAdmin); er != nil {
			err = er
			return
		} else {
			// 把 UID 挂载到上下文中国呢
			ctx.Set(ContextUidField, claims.Uid)
		}
	}
}
