// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/kataras/iris/v12"
)

var (
	ContextUidField = "uid"
)

type authQuery struct {
	Authorization string `json:"Authorization" form:"Authorization"`
}

func getToken(c iris.Context) (*string, error) {
	var query authQuery

	if err := c.ReadQuery(&query); err != nil {
		return nil, err
	}

	if len(query.Authorization) > 0 {
		return &query.Authorization, nil
	} else if len(c.GetHeader(token.AuthField)) > 0 {
		t := c.GetHeader(token.AuthField)
		return &t, nil
	} else {
		t := c.GetCookie(token.AuthField)

		if len(t) > 0 {
			return &t, nil
		}
	}

	return nil, nil
}

// Token 验证中间件
func AuthenticateNew(isAdmin bool) iris.Handler {
	return func(c iris.Context) {
		var (
			err    error
			status = schema.StatusFail
		)
		defer func() {
			if err != nil {
				_, _ = c.JSON(schema.Response{
					Status:  status,
					Message: err.Error(),
					Data:    nil,
				})
				return
			}

			c.Next()
		}()

		tokenString, err := getToken(c)

		if err != nil {
			return
		}

		if tokenString == nil {
			status = exception.InvalidToken.Code()
			err = exception.InvalidToken
			return
		}

		if claims, er := token.Parse(*tokenString, isAdmin); er != nil {
			err = er
			status = exception.InvalidToken.Code()
			return
		} else {
			// 把 UID 挂载到上下文中国呢
			c.Values().Set(ContextUidField, claims.Uid)
		}
	}
}
