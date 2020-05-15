// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin_server

import (
	"fmt"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/address"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/banner"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/help"
	loginLog "github.com/axetroy/go-server/internal/app/admin_server/controller/logger/login"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/menu"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/message"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/news"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/notification"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/report"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/role"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/system"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/user"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/dotenv"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

var AdminRouter *gin.Engine

func init() {
	if config.Common.Mode == config.ModeProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.Use(middleware.GracefulExit())

	router.Use(middleware.CORS())

	router.Static("/public", path.Join(dotenv.RootDir, "public"))

	if config.Common.Mode == config.ModeProduction {
		router.Use(gin.Logger())
	}

	router.Use(gin.Recovery())

	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, schema.Response{
			Status:  schema.StatusFail,
			Message: fmt.Sprintf("%v ", http.StatusNotFound) + http.StatusText(http.StatusNotFound),
			Data:    nil,
		})
	})

	{
		v1 := router.Group("/v1")
		v1.Use(middleware.Common)
		v1.GET("", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"ping": "pong"})
		})

		adminAuthMiddleware := middleware.Authenticate(true) // 管理员Token的中间件

		// 登陆
		v1.POST("/login", admin.LoginRouter) // 管理员登陆

		v1.Use(adminAuthMiddleware)

		v1.GET("/profile", adminAuthMiddleware, admin.GetAdminInfoRouter)    // 获取管理员自己的信息
		v1.PUT("/password", adminAuthMiddleware, admin.UpdatePasswordRouter) // 更改自己的密码

		// 管理员类
		{
			adminRouter := v1.Group("admin")
			adminRouter.POST("", admin.CreateAdminRouter)                   // 创建管理员
			adminRouter.GET("", admin.GetListRouter)                        // 获取管理员列表
			adminRouter.GET("/a/:admin_id", admin.GetAdminInfoByIdRouter)   // 获取某个管理员的信息
			adminRouter.PUT("/a/:admin_id", admin.UpdateRouter)             // 修改某个管理员的信息
			adminRouter.DELETE("/a/:admin_id", admin.DeleteAdminByIdRouter) // 修改某个管理员的信息
			adminRouter.GET("/accession", admin.GetAccessionRouter)         // 获取管理员的所有权限列表
		}

		// 用户类
		{
			userRouter := v1.Group("user")
			userRouter.GET("", user.GetListRouter)                                   // 获取会员列表
			userRouter.POST("", user.CreateUserRouter)                               // 创建会员
			userRouter.GET("/u/:user_id", user.GetProfileByAdminRouter)              // 获取单个会员的信息
			userRouter.PUT("/u/:user_id", user.UpdateProfileByAdminRouter)           // 更新会员信息
			userRouter.PUT("/u/:user_id/password", user.UpdatePasswordByAdminRouter) // 修改会员密码
		}

		// 用户角色
		{
			roleRouter := v1.Group("role")
			roleRouter.GET("", role.GetListRouter)                   // 获取角色列表
			roleRouter.POST("", role.CreateRouter)                   // 创建角色
			roleRouter.PUT("/r/:name", role.UpdateRouter)            // 修改角色
			roleRouter.DELETE("/r/:name", role.DeleteRouter)         // 删除角色
			roleRouter.GET("/r/:name", role.GetRouter)               // 获取角色详情
			roleRouter.GET("/accession", role.GetAccessionRouter)    // 获取用户的所有的权限列表
			roleRouter.GET("/u/:user_id", role.UpdateUserRoleRouter) // 用户用户的角色信息
			roleRouter.PUT("/u/:user_id", role.UpdateUserRoleRouter) // 管理员修改用户的角色
		}

		// 新闻咨询类
		{
			newsRouter := v1.Group("/news")
			newsRouter.POST("", news.CreateRouter)              // 新建新闻公告
			newsRouter.GET("", news.GetNewsListRouter)          // 获取新闻列表
			newsRouter.GET("/n/:news_id", news.GetNewsRouter)   // 获取新闻详情
			newsRouter.PUT("/n/:news_id", news.UpdateRouter)    // 更新新闻公告
			newsRouter.DELETE("/n/:news_id", news.DeleteRouter) // 删除新闻
		}

		// 系统通知
		{
			notificationRouter := v1.Group("/notification")
			notificationRouter.POST("", notification.CreateRouter)                    // 创建系统通知
			notificationRouter.GET("", notification.GetNotificationListByAdminRouter) // 获取系统通知列表
			notificationRouter.PUT("/n/:id", notification.UpdateRouter)               // 更新系统通知
			notificationRouter.DELETE("/n/:id", notification.DeleteRouter)            // 删除系统通知
			notificationRouter.GET("/n/:id", notification.GetRouter)                  // 获取单条系统通知
		}

		// 个人消息
		{
			messageRouter := v1.Group("/message")
			messageRouter.POST("", message.CreateRouter)                        // 创建个人消息
			messageRouter.GET("", message.GetMessageListByAdminRouter)          // 获取消息列表
			messageRouter.GET("/m/:message_id", message.GetAdminRouter)         // 获取个人消息
			messageRouter.PUT("/m/:message_id", message.UpdateRouter)           // 更新个人消息
			messageRouter.DELETE("/m/:message_id", message.DeleteByAdminRouter) // 删除个人消息
		}

		// 用户反馈
		{
			reportRouter := v1.Group("/report")
			reportRouter.Use(adminAuthMiddleware)
			reportRouter.GET("", report.GetListByAdminRouter)                // 获取我的反馈列表
			reportRouter.GET("/r/:report_id", report.GetReportByAdminRouter) // 获取反馈详情
			reportRouter.PUT("/r/:report_id", report.UpdateByAdminRouter)    // 更新用户反馈
		}

		// 帮助中心
		{
			helpRouter := v1.Group("help")
			helpRouter.GET("", help.GetHelpListRouter)          // 创建帮助列表
			helpRouter.POST("", help.CreateRouter)              // 创建帮助
			helpRouter.PUT("/h/:help_id", help.UpdateRouter)    // 更新帮助
			helpRouter.GET("/h/:help_id", help.GetHelpRouter)   // 获取帮助详情
			helpRouter.DELETE("/h/:help_id", help.DeleteRouter) // 删除帮助
		}

		// Banner
		{
			bannerRouter := v1.Group("banner")
			bannerRouter.GET("", banner.GetBannerListRouter)          // 获取 banner 列表
			bannerRouter.POST("", banner.CreateRouter)                // 创建 banner
			bannerRouter.PUT("/b/:banner_id", banner.UpdateRouter)    // 更新 banner
			bannerRouter.GET("/b/:banner_id", banner.GetBannerRouter) // 获取 banner 详情
			bannerRouter.DELETE("/b/:banner_id", banner.DeleteRouter) // 删除 banner
		}

		// 后台管理员菜单
		{
			menuRouter := v1.Group("menu")
			menuRouter.GET("", menu.GetListRouter)              // 获取菜单列表
			menuRouter.POST("", menu.CreateRouter)              // 创建菜单
			menuRouter.POST("/tree", menu.CreateFromTreeRouter) // 创建菜单
			menuRouter.PUT("/m/:menu_id", menu.UpdateRouter)    // 更新菜单
			menuRouter.GET("/m/:menu_id", menu.GetMenuRouter)   // 获取菜单详情
			menuRouter.DELETE("/m/:menu_id", menu.DeleteRouter) // 删除菜单
		}

		// 日志
		{
			logRouter := v1.Group("log")
			logRouter.GET("/login", loginLog.GetLoginLogsRouter)          // 获取用户的登陆日志列表
			logRouter.GET("/login/l/:log_id", loginLog.GetLoginLogRouter) // 用户单条登陆记录
		}

		// 通用类
		{
			v1.GET("/area/:area_code", address.FindAddressRouter) // 获取地区码对应的信息
			v1.GET("/area", address.AreaListRouter)               // 获取地址选择列表
		}

		v1.GET("/system", system.GetSystemInfoRouter) // 获取系统相关信息
	}

	AdminRouter = router
}
