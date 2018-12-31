package router

import (
	"github.com/gin-gonic/gin"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/controller/downloader"
	"github.com/axetroy/go-server/controller/email"
	"github.com/axetroy/go-server/controller/finance"
	"github.com/axetroy/go-server/controller/invite"
	"github.com/axetroy/go-server/controller/news"
	"github.com/axetroy/go-server/controller/resource"
	"github.com/axetroy/go-server/controller/static"
	"github.com/axetroy/go-server/controller/transfer"
	"github.com/axetroy/go-server/controller/uploader"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/controller/wallet"
	"github.com/axetroy/go-server/middleware"
	"net/http"
)

var Router *gin.Engine

func init() {
	router := gin.Default()

	router.Use(gin.Logger())

	// Recovery middleware recovers from any panics and writes a 500 if there was one.
	router.Use(gin.Recovery())

	router.GET("/ping", func(context *gin.Context) {
		context.String(http.StatusOK, "pong")
	})

	// Simple group: v1
	v1 := router.Group("/v1")

	v1.Use(func(context *gin.Context) {
		header := context.Writer.Header()
		// alone dns prefect
		header.Set("X-DNS-Prefetch-Control", "on")
		// IE No Open
		header.Set("X-Download-Options", "noopen")
		// not cache
		header.Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
		header.Set("Expires", "max-age=0")
		// Content Security Policy
		header.Set("Content-Security-Policy", "default-src 'self'")
		// xss protect
		// it will caught some problems is old IE
		header.Set("X-XSS-Protection", "1; mode=block")
		// Referrer Policy
		header.Set("Referrer-Header", "no-referrer")
		// cros frame, allow same origin
		header.Set("X-Frame-Options", "SAMEORIGIN")
		// HSTS
		header.Set("Strict-Transport-Security", "max-age=5184000;includeSubDomains")
		// no sniff
		header.Set("X-Content-Type-Options", "nosniff")
	})

	{
		v1.GET("", func(context *gin.Context) {
			context.JSON(http.StatusOK, gin.H{
				"ping": "pong",
			})
		})
		// 认证类
		authRouter := v1.Group("/auth")
		{
			authRouter.POST("/signup", auth.SignUp)
			authRouter.POST("/signin", auth.Signin)
			authRouter.POST("/activation", auth.Activation)
			authRouter.PUT("/password/reset", auth.ResetPassword)
		}

		// 用户类
		userRouter := v1.Group("/user")
		{
			userRouter.Use(middleware.Authenticate())
			userRouter.GET("/signout", user.Signout)
			userRouter.GET("/profile", user.GetProfile)
			userRouter.PUT("/profile", user.UpdateProfile)
			userRouter.PUT("/password/update", user.UpdatePassword)
			userRouter.PUT("/trade_password/set", user.SetPayPassword)
			userRouter.PUT("/trade_password/update", user.UpdatePayPassword)
			// TODO: 上传头像
			userRouter.POST("/avatar", user.UpdatePayPassword)
			// 邀请人列表
			userRouter.GET("/invite", invite.GetMyInviteList)
		}

		// 钱包类
		walletRouter := v1.Group("/wallet")
		{
			walletRouter.Use(middleware.Authenticate())
			// 获取所有的钱包信息
			walletRouter.GET("/map", wallet.GetWallets)
			walletRouter.GET("/currency/:currency", wallet.GetWallet)
			// 转账相关
			walletRouter.GET("/transfer/history", transfer.GetHistory)
			walletRouter.GET("/transfer/detail/:id", transfer.GetDetail)
			walletRouter.POST("/transfer", middleware.AuthPayPassword, transfer.To)
		}

		// 财务日志
		financeRouter := v1.Group("/finance")
		{
			financeRouter.Use(middleware.Authenticate())
			financeRouter.GET("/history", finance.GetHistory) // TODO: 获取我的财务日志
		}

		// 新闻咨询类
		newsRouter := v1.Group("/news")
		{
			// TODO: 写新闻咨询类
			newsRouter.Use(middleware.Authenticate())
			newsRouter.POST("/", news.Create)
			newsRouter.GET("/list", news.GetNewsList)
			newsRouter.GET("/detail/:id", news.GetNews)
			newsRouter.PUT("/update/:id", news.Update)
		}

		// 系统通知
		notificationRouter := v1.Group("/notification")
		{
			// TODO: 写通知类
			notificationRouter.GET("/")
			notificationRouter.GET("/:id")
		}

		// 用户的个人通知
		messageRouter := v1.Group("/message")
		{
			// TODO: 写个人通知
			messageRouter.GET("/")
			messageRouter.GET("/:id")
		}

		// 通用类
		{
			// 邮件服务
			v1.POST("/email/send/activation", email.SendActivationEmail)
			v1.POST("/email/send/reset_password", email.SendResetPasswordEmail)

			// 文件上传 (需要验证token)
			uploadRouter := v1.Group("/upload")
			{
				uploadRouter.Use(middleware.Authenticate())
				uploadRouter.POST("/file", uploader.File)
				uploadRouter.POST("/image", uploader.Image)
			}
			// 单纯获取资源文本
			v1.GET("/resource/file/:filename", resource.File)
			v1.GET("/resource/image/:filename", resource.Image)
			v1.GET("/resource/thumbnail/:filename", resource.Thumbnail)
			// 下载资源
			v1.GET("/download/file/:filename", downloader.File)
			v1.GET("/download/image/:filename", downloader.Image)
			v1.GET("/download/thumbnail/:filename", downloader.Thumbnail)
			// 公共资源目录
			v1.GET("/public/:filename", static.Get)
		}
	}

	Router = router
}
