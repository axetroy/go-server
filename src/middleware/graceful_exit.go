// Copyright 2019 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"github.com/axetroy/go-server/src/config"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 优雅退出中间件
// 再接收到退出指令之后，则 HTTP 服务不再接收新的请求
func GracefulExit() gin.HandlerFunc {
	return func(context *gin.Context) {
		if config.Common.Exiting {
			err := exception.SystemMaintenance
			context.JSON(http.StatusOK, schema.Response{
				Status:  err.Code(),
				Message: err.Error(),
				Data:    nil,
			})
			context.Abort()
		}
	}
}
