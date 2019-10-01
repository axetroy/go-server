// Copyright 2019 Axetroy. All rights reserved. MIT license.
package controller

import (
	"github.com/axetroy/go-server/src/middleware"
	"github.com/gin-gonic/gin"
)

// 控制器的上下文
type Context struct {
	*gin.Context
	Uid       string `json:"uid"`        // 操作人的用户 ID
	UserAgent string `json:"user_agent"` // 用户代理
	Ip        string `json:"ip"`         // IP地址
}

func NewContext(c *gin.Context) Context {
	return Context{
		Uid:       c.GetString(middleware.ContextUidField),
		UserAgent: c.GetHeader("user-agent"),
		Ip:        c.ClientIP(),
	}
}
