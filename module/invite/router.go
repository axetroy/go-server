package invite

import (
	"github.com/axetroy/go-server/middleware"
	"github.com/gin-gonic/gin"
)

func Route(r *gin.RouterGroup) *gin.RouterGroup {
	inviteRouter := r.Group("/invite")

	inviteRouter.Use(middleware.Authenticate(false))
	inviteRouter.GET("", GetListRouter)          // 获取我已邀请的列表
	inviteRouter.GET("/i/:invite_id", GetRouter) // 获取单条邀请记录详情

	return r
}
