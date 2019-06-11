// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address

import "github.com/gin-gonic/gin"

func Route(r *gin.RouterGroup) *gin.RouterGroup {
	addressRouter := r.Group("/address")

	addressRouter.GET("", GetListRouter)                 // 获取地址列表
	addressRouter.POST("", CreateRouter)                 // 添加收货地址
	addressRouter.PUT("/a/:address_id", UpdateRouter)    // 更新收货地址
	addressRouter.DELETE("/a/:address_id", DeleteRouter) // 删除收货地址
	addressRouter.GET("/a/:address_id", GetDetailRouter) // 获取地址详情
	addressRouter.GET("/default", GetDefaultRouter)      // 获取默认地址

	return addressRouter
}
