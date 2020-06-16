// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin_server

import (
	"fmt"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/area"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/banner"
	Configuration "github.com/axetroy/go-server/internal/app/admin_server/controller/config"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/help"
	loginLog "github.com/axetroy/go-server/internal/app/admin_server/controller/logger/login"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/menu"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/message"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/news"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/notification"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/push"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/report"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/role"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/system"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/user"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"net/http"
)

var AdminRouter *iris.Application

func init() {
	app := iris.New()

	app.OnAnyErrorCode(router.Handler(func(c router.Context) {
		code := c.GetStatusCode()

		fmt.Println(c.Request().Method, c.Request().URL.Path, code)

		c.StatusCode(code)

		c.JSON(fmt.Errorf("%d %s", code, http.StatusText(code)), nil, nil)
	}))

	{
		v1 := app.Party("/v1").AllowMethods(iris.MethodOptions)
		v1.Use(recover.New())
		v1.Use(middleware.Common())
		v1.Use(middleware.CORS())
		v1.Use(middleware.RateLimit(30))

		if config.Common.Mode != "production" {
			v1.Use(logger.New())
			v1.Use(middleware.Ip())
		}

		{
			v1.Get("", router.Handler(func(c router.Context) {
				c.JSON(nil, map[string]string{"ping": "tong"}, nil)
			}))
		}

		adminAuthMiddleware := middleware.AuthenticateNew(true) // 管理员Token的中间件

		// 登陆
		v1.Post("/login", admin.LoginRouter) // 管理员登陆

		v1.Use(adminAuthMiddleware)

		v1.Get("/profile", adminAuthMiddleware, admin.GetAdminInfoRouter)    // 获取管理员自己的信息
		v1.Put("/password", adminAuthMiddleware, admin.UpdatePasswordRouter) // 更改自己的密码

		// 管理员类
		{
			adminRouter := v1.Party("/admin")
			adminRouter.Post("", admin.CreateAdminRouter)                  // 创建管理员
			adminRouter.Get("", admin.GetListRouter)                       // 获取管理员列表
			adminRouter.Get("/accession", admin.GetAccessionRouter)        // 获取管理员的所有权限列表
			adminRouter.Get("/{admin_id}", admin.GetAdminInfoByIdRouter)   // 获取某个管理员的信息
			adminRouter.Put("/{admin_id}", admin.UpdateRouter)             // 修改某个管理员的信息
			adminRouter.Delete("/{admin_id}", admin.DeleteAdminByIdRouter) // 修改某个管理员的信息
		}

		// 用户类
		{
			userRouter := v1.Party("/user")
			userRouter.Get("", user.GetListRouter)                                  // 获取会员列表
			userRouter.Post("", user.CreateUserRouter)                              // 创建会员
			userRouter.Get("/{user_id}", user.GetProfileByAdminRouter)              // 获取单个会员的信息
			userRouter.Put("/{user_id}/password", user.UpdatePasswordByAdminRouter) // 修改会员密码
			userRouter.Put("/{user_id}", user.UpdateProfileByAdminRouter)           // 更新会员信息
			userRouter.Put("/{user_id}/role", role.UpdateUserRoleRouter)            // 修改用户的角色
		}

		// 用户角色
		{
			roleRouter := v1.Party("/role")
			roleRouter.Get("", role.GetListRouter)                // 获取角色列表
			roleRouter.Post("", role.CreateRouter)                // 创建角色
			roleRouter.Put("/{name}", role.UpdateRouter)          // 修改角色
			roleRouter.Delete("/{name}", role.DeleteRouter)       // 删除角色
			roleRouter.Get("/{name}", role.GetRouter)             // 获取角色详情
			roleRouter.Get("/accession", role.GetAccessionRouter) // 获取用户的所有的权限列表
		}

		// 新闻咨询类
		{
			newsRouter := v1.Party("/news")
			newsRouter.Post("", news.CreateRouter)             // 新建新闻公告
			newsRouter.Get("", news.GetNewsListRouter)         // 获取新闻列表
			newsRouter.Get("/{news_id}", news.GetNewsRouter)   // 获取新闻详情
			newsRouter.Put("/{news_id}", news.UpdateRouter)    // 更新新闻公告
			newsRouter.Delete("/{news_id}", news.DeleteRouter) // 删除新闻
		}

		// 系统通知
		{
			notificationRouter := v1.Party("/notification")
			notificationRouter.Post("", notification.CreateRouter)                    // 创建系统通知
			notificationRouter.Get("", notification.GetNotificationListByAdminRouter) // 获取系统通知列表
			notificationRouter.Put("/{id}", notification.UpdateRouter)                // 更新系统通知
			notificationRouter.Delete("/{id}", notification.DeleteRouter)             // 删除系统通知
			notificationRouter.Get("/{id}", notification.GetRouter)                   // 获取单条系统通知
		}

		// 个人消息
		{
			messageRouter := v1.Party("/message")
			messageRouter.Post("", message.CreateRouter)                       // 创建个人消息
			messageRouter.Get("", message.GetMessageListByAdminRouter)         // 获取消息列表
			messageRouter.Get("/{message_id}", message.GetAdminRouter)         // 获取个人消息
			messageRouter.Put("/{message_id}", message.UpdateRouter)           // 更新个人消息
			messageRouter.Delete("/{message_id}", message.DeleteByAdminRouter) // 删除个人消息
		}

		// 用户反馈
		{
			reportRouter := v1.Party("/report")
			reportRouter.Use(adminAuthMiddleware)
			reportRouter.Get("/type", report.GetTypesRouter)                // 获取类型列表
			reportRouter.Get("", report.GetListByAdminRouter)               // 获取我的反馈列表
			reportRouter.Get("/{report_id}", report.GetReportByAdminRouter) // 获取反馈详情
			reportRouter.Put("/{report_id}", report.UpdateByAdminRouter)    // 更新用户反馈
		}

		// 帮助中心
		{
			helpRouter := v1.Party("/help")
			helpRouter.Get("", help.GetHelpListRouter)         // 创建帮助列表
			helpRouter.Post("", help.CreateRouter)             // 创建帮助
			helpRouter.Put("/{help_id}", help.UpdateRouter)    // 更新帮助
			helpRouter.Get("/{help_id}", help.GetHelpRouter)   // 获取帮助详情
			helpRouter.Delete("/{help_id}", help.DeleteRouter) // 删除帮助
		}

		// Banner
		{
			bannerRouter := v1.Party("/banner")
			bannerRouter.Get("", banner.GetBannerListRouter)         // 获取 banner 列表
			bannerRouter.Post("", banner.CreateRouter)               // 创建 banner
			bannerRouter.Put("/{banner_id}", banner.UpdateRouter)    // 更新 banner
			bannerRouter.Get("/{banner_id}", banner.GetBannerRouter) // 获取 banner 详情
			bannerRouter.Delete("/{banner_id}", banner.DeleteRouter) // 删除 banner
		}

		// 后台管理员菜单
		{
			menuRouter := v1.Party("/menu")
			menuRouter.Get("", menu.GetListRouter)             // 获取菜单列表
			menuRouter.Post("", menu.CreateRouter)             // 创建菜单
			menuRouter.Get("/tree", menu.CreateFromTreeRouter) // 创建菜单
			menuRouter.Put("/{menu_id}", menu.UpdateRouter)    // 更新菜单
			menuRouter.Get("/{menu_id}", menu.GetMenuRouter)   // 获取菜单详情
			menuRouter.Delete("/{menu_id}", menu.DeleteRouter) // 删除菜单
		}

		// 日志
		{
			logRouter := v1.Party("/log")
			logRouter.Get("/login", loginLog.GetLoginLogsRouter)         // 获取用户的登陆日志列表
			logRouter.Get("/login/{log_id}", loginLog.GetLoginLogRouter) // 用户单条登陆记录
		}

		// 配置项
		{
			configRouter := v1.Party("/config")
			configRouter.Get("/", Configuration.GetListRouter)             // 获取配置列表
			configRouter.Get("/name", Configuration.GetNameRouter)         // 获取配置名称列表 (查看有哪些配置)
			configRouter.Get("/:config_name", Configuration.GetRouter)     // 获取指定的配置
			configRouter.Post("/:config_name", Configuration.CreateRouter) // 创建指定的配置
			configRouter.Put("/:config_name", Configuration.UpdateRouter)  // 更新指定的配置
		}

		// 推送
		{
			configRouter := v1.Party("/push")
			configRouter.Get("/notification", push.CreateNotificationRouter)  // TODO: 获取推送列表
			configRouter.Post("/notification", push.CreateNotificationRouter) // 生成一个推送到指定用户
		}

		// 地区接口
		{
			areaRouter := v1.Party("/area")
			areaRouter.Get("/provinces", area.GetProvinces)      // 获取省份
			areaRouter.Get("/{code}", area.GetDetail)            // 获取地区码的详情
			areaRouter.Get("/{code}/children", area.GetChildren) // 获取地区下的子地区
			areaRouter.Get("", area.GetArea)                     // 获取所有地区
		}

		v1.Get("/system", system.GetSystemInfoRouter) // 获取系统相关信息
	}

	_ = app.Build()

	AdminRouter = app
}
