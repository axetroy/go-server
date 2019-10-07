// Copyright 2019 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ContextUidField = "uid"
)

// Token 验证中间件
func Authenticate(isAdmin bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		var (
			err         error
			tokenString string
			status      = schema.StatusFail
		)
		defer func() {
			if err != nil {
				c.JSON(http.StatusOK, schema.Response{
					Status:  status,
					Message: err.Error(),
					Data:    nil,
				})
				c.Abort()
			}
		}()

		if s, isExist := c.GetQuery(token.AuthField); isExist == true {
			tokenString = s
			return
		} else {
			tokenString = c.GetHeader(token.AuthField)

			if len(tokenString) == 0 {
				if s, er := c.Cookie(token.AuthField); er != nil {
					err = exception.InvalidToken
					status = exception.InvalidToken.Code()
					return
				} else {
					tokenString = s
				}
			}
		}

		if claims, er := token.Parse(tokenString, isAdmin); er != nil {
			err = er
			status = exception.InvalidToken.Code()
			return
		} else {
			// 把 UID 挂载到上下文中国呢
			c.Set(ContextUidField, claims.Uid)
		}
	}
}
