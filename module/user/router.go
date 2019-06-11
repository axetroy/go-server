// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/rbac"
	"github.com/axetroy/go-server/rbac/accession"
	"github.com/gin-gonic/gin"
)

func Route(router *gin.RouterGroup) *gin.RouterGroup {
	userRouter := router.Group("/user")

	userRouter.Use(middleware.Authenticate(false))
	userRouter.GET("/signout", SignOut)                                                                      // 用户登出
	userRouter.GET("/profile", GetProfileRouter)                                                             // 获取用户详细信息
	userRouter.PUT("/profile", rbac.Require(*accession.ProfileUpdate), UpdateProfileRouter)                  // 更新用户资料
	userRouter.PUT("/password", rbac.Require(*accession.PasswordUpdate), UpdatePasswordRouter)               // 更新登陆密码
	userRouter.POST("/password2", rbac.Require(*accession.Password2Set), SetPayPasswordRouter)               // 设置交易密码
	userRouter.PUT("/password2", rbac.Require(*accession.Password2Update), UpdatePayPasswordRouter)          // 更新交易密码
	userRouter.PUT("/password2/reset", rbac.Require(*accession.Password2Reset), ResetPayPasswordRouter)      // 重置交易密码
	userRouter.POST("/password2/reset", rbac.Require(*accession.Password2Reset), SendResetPayPasswordRouter) // 发送重置交易密码的邮件/短信
	userRouter.POST("/avatar", UploadAvatarRouter)                                                           // 上传用户头像

	return router
}
