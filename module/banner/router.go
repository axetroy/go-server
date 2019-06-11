// Copyright 2019 Axetroy. All rights reserved. MIT license.

package banner

import "github.com/gin-gonic/gin"

func Route(r *gin.RouterGroup) *gin.RouterGroup {
	bannerRouter := r.Group("banner")

	bannerRouter.GET("", GetListRouter)                // 获取 banner 列表
	bannerRouter.GET("/b/:banner_id", GetBannerRouter) // 获取 banner 详情

	return bannerRouter
}
