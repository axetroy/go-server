// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import "github.com/gin-gonic/gin"

func Route(r *gin.RouterGroup) *gin.RouterGroup {
	authRouter := r.Group("auth")

	authRouter.POST("/signup", SignUpRouter)               // 注册账号
	authRouter.POST("/signin", SignInRouter)               // 登陆账号
	authRouter.POST("/activation", ActivationRouter)       // 激活账号
	authRouter.PUT("/password/reset", ResetPasswordRouter) // 密码重置

	return authRouter
}
