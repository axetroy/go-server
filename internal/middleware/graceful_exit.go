// Copyright 2019 Axetroy. All rights reserved. MIT license.
package middleware

import (
	"github.com/axetroy/go-server/internal/config"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 优雅退出中间件
// 再接收到退出指令之后，则 HTTP 服务不再接收新的请求
func GracefulExit() gin.HandlerFunc {
	return func(c *gin.Context) {
		if config.Common.Exiting {
			err := exception.SystemMaintenance
			c.JSON(http.StatusOK, schema.Response{
				Status:  err.Code(),
				Message: err.Error(),
				Data:    nil,
			})
			c.Abort()
		}
	}
}
