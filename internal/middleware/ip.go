// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"github.com/kataras/iris/v12"
)

//func Ip() iris.Handler {
//	return router.Handler(func(c router.Context) {
//		c.Header("X-Client-Ip", c.ClientIP())
//
//		fmt.Println("当前 IP", c.ClientIP())
//
//		c.Next()
//	})
//}

func Ip() iris.Handler {
	return func(c iris.Context) {
		c.Header("X-Client-Ip", c.RemoteAddr())
		c.Next()
	}
}
