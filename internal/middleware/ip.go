// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/kataras/iris/v12"
)

func Ip() iris.Handler {
	return router.Handler(func(c router.Context) {
		c.Header("X-Client-Ip", c.ClientIP())

		c.Next()
	})
}
