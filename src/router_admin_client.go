package src

import (
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/banner"
	"github.com/axetroy/go-server/src/controller/message"
	"github.com/axetroy/go-server/src/controller/news"
	"github.com/axetroy/go-server/src/controller/notification"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

var RouterAdminClient *gin.Engine

func init() {
	gin.SetMode(gin.ReleaseMode)
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
			context.JSON(http.StatusOK, gin.H{"ping": "pong"})
		})

		adminAuthMiddleware := middleware.Authenticate(true) // 管理员Token的中间件

		// 登陆
		v1.POST("/login", admin.LoginRouter)

		v1.Use(adminAuthMiddleware)

		// 管理员类
		adRouter := v1.Group("admin")
		{
			adRouter.POST("", admin.CreateAdminRouter)         // 创建管理员
			adRouter.GET("", admin.GetAdminInfoRouter)         // TODO: 获取管理员列表
			adRouter.GET("/profile", admin.GetAdminInfoRouter) // 获取管理员信息

			u := v1.Group("u")
			{
				u2 := u.Group(":admin_id")
				{
					u2.GET("", admin.GetAdminInfoRouter)          // TODO: 获取某个管理员的信息
					u2.PUT("", admin.GetAdminInfoRouter)          // TODO: 修改管理员信息
					u2.DELETE("", admin.GetAdminInfoRouter)       // TODO: 删除管理员
					u2.PUT("/password", admin.GetAdminInfoRouter) // TODO: 修改管理员密码
				}
			}
		}

		// 用户类
		userRouter := v1.Group("user")
		{
			u := userRouter.Group("u")
			{
				u.GET("/:user_id", admin.GetAdminInfoRouter)    // TODO: 获取单个会员的信息
				u.PUT("/:user_id", admin.GetAdminInfoRouter)    // TODO: 更新会员信息
				u.DELETE("/:user_id", admin.GetAdminInfoRouter) // TODO: 删除会员信息
			}
			userRouter.GET("/", admin.GetAdminInfoRouter)         // TODO: 获取会员列表
			userRouter.PUT("/password", admin.GetAdminInfoRouter) // TODO: 修改会员密码
		}

		// 新闻咨询类
		newsRouter := v1.Group("/news")
		{
			newsRouter.POST("", news.CreateRouter)
			newsRouter.GET("", news.UpdateRouter) // TODO: 获取新闻列表

			n := newsRouter.Group("n")
			{
				n.PUT("/:news_id", news.UpdateRouter)
				n.DELETE("/:news_id", news.UpdateRouter) // TODO: 删除新闻
				n.GET("/:news_id", news.UpdateRouter)    // TODO: 获取新闻详情
			}
		}

		// 系统通知
		notificationRouter := v1.Group("/notification")
		{
			notificationRouter.POST("", notification.CreateRouter)
			notificationRouter.GET("", notification.DeleteRouter) // TODO: 获取系统通知列表

			n := notificationRouter.Group("n")
			{
				n.PUT("/:id", notification.UpdateRouter)
				n.DELETE("/:id", notification.DeleteRouter)
				n.GET("/:id", notification.DeleteRouter) // TODO: 获取单条系统通知
			}

		}

		// 个人消息
		messageRouter := v1.Group("/message")
		{
			messageRouter.GET("/", message.DeleteByAdminRouter) // TODO: 获取个人消息列表
			messageRouter.POST("/", message.CreateRouter)

			m := messageRouter.Group("/m/:message_id")

			{
				m.PUT("/", message.UpdateRouter)
				m.GET("/", message.DeleteByAdminRouter) // TODO: 获取个人消息
				m.DELETE("/", message.DeleteByAdminRouter)
			}
		}

		// Banner
		bannerRouter := v1.Group("banner")
		{
			bannerRouter.GET("", banner.UpdateRouter) // TODO: 获取 banner 列表
			bannerRouter.POST("", banner.CreateRouter)

			b := bannerRouter.Group("b")
			{
				b.PUT("/:banner_id", banner.UpdateRouter)
				b.GET("/:banner_id", banner.UpdateRouter)    // TODO: 获取 banner 详情
				b.DELETE("/:banner_id", banner.UpdateRouter) // TODO: 删除 banner
			}
		}
	}

	RouterAdminClient = router
}
