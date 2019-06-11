// Copyright 2019 Axetroy. All rights reserved. MIT license.

package email

import "github.com/gin-gonic/gin"

func Route(r *gin.RouterGroup) *gin.RouterGroup {

	r.POST("/email/send/activation", SendActivationEmailRouter)        // 发送激活邮件
	r.POST("/email/send/password/reset", SendResetPasswordEmailRouter) // 发送密码重置邮件

	return r
}
