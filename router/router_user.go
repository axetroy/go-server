// Copyright 2019 Axetroy. All rights reserved. MIT license.
package router

import (
	"fmt"
	"github.com/axetroy/go-server/config"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/address"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/banner"
	"github.com/axetroy/go-server/module/downloader"
	"github.com/axetroy/go-server/module/email"
	"github.com/axetroy/go-server/module/finance"
	"github.com/axetroy/go-server/module/invite"
	"github.com/axetroy/go-server/module/message"
	"github.com/axetroy/go-server/module/news"
	"github.com/axetroy/go-server/module/notification"
	"github.com/axetroy/go-server/module/oauth2"
	"github.com/axetroy/go-server/module/report"
	"github.com/axetroy/go-server/module/resource"
	"github.com/axetroy/go-server/module/transfer"
	"github.com/axetroy/go-server/module/uploader"
	"github.com/axetroy/go-server/module/user"
	"github.com/axetroy/go-server/module/wallet"
	"github.com/axetroy/go-server/rbac"
	"github.com/axetroy/go-server/rbac/accession"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/dotenv"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

var UserRouter *gin.Engine

func init() {
	if config.Common.Mode == config.ModeProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()
	router.Use(middleware.CORS())
	router.Static("/public", path.Join(dotenv.RootDir, "public"))

	if config.Common.Mode != config.ModeProduction {
		router.Use(gin.Logger())
	}
	router.Use(gin.Recovery())

	router.NoRoute(func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, schema.Response{
			Status:  schema.StatusFail,
			Message: fmt.Sprintf("%v ", http.StatusNotFound) + http.StatusText(http.StatusNotFound),
			Data:    nil,
		})
	})

	{
		v1 := router.Group("/v1")
		v1.Use(middleware.Common)

		v1.GET("", func(ctx *gin.Context) {
			ctx.JSON(http.StatusOK, gin.H{"ping": "pong"})
		})

		userAuthMiddleware := middleware.Authenticate(false) // 用户Token的中间件

		// 认证类
		auth.Route(v1) // 验证模块

		// oAuth2 认证
		{
			oAuthRouter := v1.Group("/oauth2")
			oAuthRouter.GET("/google", oauth2.GoogleLoginRouter)             // 用 Google 登陆
			oAuthRouter.GET("/google_callback", oauth2.GoogleCallbackRouter) // Google 认证完成后跳转到这里，用户不应该访问这个地址
		}

		user.Route(v1)    // 用户类模块
		invite.Route(v1)  // 用户邀请模块
		address.Route(v1) // 收货地址模块

		// 钱包类
		{
			walletRouter := v1.Group("/wallet")
			walletRouter.Use(userAuthMiddleware)
			walletRouter.GET("", wallet.GetWalletsRouter)            // 获取所有钱包列表
			walletRouter.GET("/w/:currency", wallet.GetWalletRouter) // 获取单个钱包的详细信息
		}

		{
			transferRouter := v1.Group("/transfer")
			transferRouter.Use(userAuthMiddleware)
			transferRouter.GET("", transfer.GetHistoryRouter)                                                           // 获取我的转账记录
			transferRouter.POST("", rbac.Require(*accession.DoTransfer), middleware.AuthPayPassword, transfer.ToRouter) // 转账给某人
			transferRouter.GET("/t/:transfer_id", transfer.GetDetailRouter)                                             // 获取单条转账详情
		}

		// 财务日志
		{
			financeRouter := v1.Group("/finance")
			financeRouter.Use(userAuthMiddleware)
			financeRouter.GET("/history", finance.GetHistory) // TODO: 获取我的财务日志
		}

		// 新闻咨询类
		{
			newsRouter := v1.Group("/news")
			newsRouter.GET("", news.GetListRouter)       // 获取新闻公告列表
			newsRouter.GET("/n/:id", news.GetNewsRouter) // 获取单个新闻公告详情
		}

		// 系统通知
		{
			notificationRouter := v1.Group("/notification")
			notificationRouter.Use(userAuthMiddleware)
			notificationRouter.GET("", notification.GetListUserRouter)     // 获取系统通知列表
			notificationRouter.GET("/n/:id", notification.GetRouter)       // 获取某一条系统通知详情
			notificationRouter.PUT("/n/:id/read", notification.ReadRouter) // 标记通知为已读
		}

		// 用户的个人消息, 个人消息是可以删除的
		{
			messageRouter := v1.Group("/message")
			messageRouter.Use(userAuthMiddleware)
			messageRouter.GET("", message.GetListRouter)                       // 获取我的消息列表
			messageRouter.GET("/m/:message_id", message.GetRouter)             // 获取单个消息详情
			messageRouter.PUT("/m/:message_id/read", message.ReadRouter)       // 标记消息为已读
			messageRouter.DELETE("/m/:message_id", message.DeleteByUserRouter) // 删除消息
		}

		// 用户反馈
		{
			reportRouter := v1.Group("/report")
			reportRouter.Use(userAuthMiddleware)
			reportRouter.GET("", report.GetListRouter)                // 获取我的反馈列表
			reportRouter.POST("", report.CreateRouter)                // 添加一条反馈
			reportRouter.GET("/r/:report_id", report.GetReportRouter) // 获取反馈详情
			reportRouter.PUT("/r/:report_id", report.UpdateRouter)    // 更新这条反馈信息
		}

		banner.Route(v1)

		// 通用类
		{
			// 邮件服务
			email.Route(v1)

			// 文件上传
			v1.POST("/upload/file", uploader.File)      // 上传文件
			v1.POST("/upload/image", uploader.Image)    // 上传图片
			v1.GET("/upload/example", uploader.Example) // 上传文件的 example
			// 单纯获取资源文本
			v1.GET("/resource/file/:filename", resource.File)           // 获取文件纯文本
			v1.GET("/resource/image/:filename", resource.Image)         // 获取图片纯文本
			v1.GET("/resource/thumbnail/:filename", resource.Thumbnail) // 获取缩略图纯文本
			// 下载资源
			downloader.Route(v1)
			// 公共资源目录
			v1.GET("/avatar/:filename", user.GetAvatarRouter) // 获取用户头像

			v1.GET("/area", address.AreaListRouter) // 获取地址选择列表
		}

	}

	UserRouter = router
}
