// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin

import "github.com/gin-gonic/gin"

func Route(r *gin.RouterGroup) *gin.RouterGroup {
	adminRouter := r.Group("admin")

	adminRouter.POST("", CreateAdminRouter)                 // 创建管理员
	adminRouter.GET("", GetListRouter)                      // 获取管理员列表
	adminRouter.GET("/a/:admin_id", GetAdminInfoByIdRouter) // 获取某个管理员的信息
	adminRouter.PUT("/a/:admin_id", UpdateRouter)           // 修改某个管理员的信息

	return adminRouter
}
