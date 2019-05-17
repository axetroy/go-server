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

	v1.Use(middleware.Common)

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

			adRouter.GET("/u/:admin_id", admin.GetAdminInfoRouter)          // TODO: 获取某个管理员的信息
			adRouter.PUT("/u/:admin_id", admin.GetAdminInfoRouter)          // TODO: 修改管理员信息
			adRouter.DELETE("/u/:admin_id", admin.GetAdminInfoRouter)       // TODO: 删除管理员
			adRouter.PUT("/u/:admin_id/password", admin.GetAdminInfoRouter) // TODO: 修改管理员密码
		}

		// 用户类
		userRouter := v1.Group("user")
		{
			userRouter.GET("/", admin.GetAdminInfoRouter)              // TODO: 获取会员列表
			userRouter.PUT("/password", admin.GetAdminInfoRouter)      // TODO: 修改会员密码
			userRouter.GET("/u/:user_id", admin.GetAdminInfoRouter)    // TODO: 获取单个会员的信息
			userRouter.PUT("/u/:user_id", admin.GetAdminInfoRouter)    // TODO: 更新会员信息
			userRouter.DELETE("/u/:user_id", admin.GetAdminInfoRouter) // TODO: 删除会员信息
		}

		// 新闻咨询类
		newsRouter := v1.Group("/news")
		{
			newsRouter.POST("", news.CreateRouter)              // 新建新闻公告
			newsRouter.GET("", news.GetListRouter)              // 获取新闻列表
			newsRouter.GET("/n/:news_id", news.GetNewsRouter)   // 获取新闻详情
			newsRouter.PUT("/n/:news_id", news.UpdateRouter)    // 更新新闻公告
			newsRouter.DELETE("/n/:news_id", news.DeleteRouter) // 删除新闻
		}

		// 系统通知
		notificationRouter := v1.Group("/notification")
		{
			notificationRouter.POST("", notification.CreateRouter)         // 创建系统通知
			notificationRouter.GET("", notification.DeleteRouter)          // TODO: 获取系统通知列表
			notificationRouter.PUT("/n/:id", notification.UpdateRouter)    // 更新系统通知
			notificationRouter.DELETE("/n/:id", notification.DeleteRouter) // 删除系统通知
			notificationRouter.GET("/n/:id", notification.GetRouter)       // 获取单条系统通知
		}

		// 个人消息
		messageRouter := v1.Group("/message")
		{
			messageRouter.POST("", message.CreateRouter)                        // 创建个人消息
			messageRouter.GET("", message.GetListAdminRouter)                   // 获取消息列表
			messageRouter.GET("/m/:message_id", message.GetAdminRouter)         // 获取个人消息
			messageRouter.PUT("/m/:message_id", message.UpdateRouter)           // 更新个人消息
			messageRouter.DELETE("/m/:message_id", message.DeleteByAdminRouter) // 删除个人消息
		}

		// Banner
		bannerRouter := v1.Group("banner")
		{
			bannerRouter.GET("", banner.GetListRouter)                // 获取 banner 列表
			bannerRouter.POST("", banner.CreateRouter)                // 创建 banner
			bannerRouter.PUT("/b/:banner_id", banner.UpdateRouter)    // 更新 banner
			bannerRouter.GET("/b/:banner_id", banner.GetBannerRouter) // 获取 banner 详情
			bannerRouter.DELETE("/b/:banner_id", banner.DeleteRouter) // 删除 banner
		}
	}

	RouterAdminClient = router
}
