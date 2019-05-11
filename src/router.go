package src

import (
	"github.com/axetroy/go-server/src/controller/address"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/banner"
	"github.com/axetroy/go-server/src/controller/downloader"
	"github.com/axetroy/go-server/src/controller/email"
	"github.com/axetroy/go-server/src/controller/finance"
	"github.com/axetroy/go-server/src/controller/invite"
	"github.com/axetroy/go-server/src/controller/message"
	"github.com/axetroy/go-server/src/controller/news"
	"github.com/axetroy/go-server/src/controller/notification"
	"github.com/axetroy/go-server/src/controller/oauth2"
	"github.com/axetroy/go-server/src/controller/resource"
	"github.com/axetroy/go-server/src/controller/static"
	"github.com/axetroy/go-server/src/controller/transfer"
	"github.com/axetroy/go-server/src/controller/uploader"
	"github.com/axetroy/go-server/src/controller/user"
	"github.com/axetroy/go-server/src/controller/wallet"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

var Router *gin.Engine

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

		userAuthMiddleware := middleware.Authenticate(false) // 用户Token的中间件
		adminAuthMiddleware := middleware.Authenticate(true) // 管理员Token的中间件

		// 管理员所有接口
		// TODO: 管理员接口应该和用户接口分离
		adminRouter := v1.Group("/admin")
		{
			// 登陆
			adminRouter.POST("/login", admin.LoginRouter)

			// 管理员类
			adRouter := adminRouter.Group("admin")
			{
				adRouter.Use(adminAuthMiddleware)
				adRouter.POST("/create", admin.CreateAdminRouter)
				adRouter.POST("/profile", admin.GetAdminInfoRouter)
			}

			// 新闻咨询类
			newsRouter := adminRouter.Group("/news")
			{
				newsRouter.Use(adminAuthMiddleware)
				newsRouter.POST("/create", news.CreateRouter)
				newsRouter.PUT("/update/:news_id", news.UpdateRouter)
			}

			// 系统通知
			notificationRouter := adminRouter.Group("/notification")
			{
				notificationRouter.POST("/create", notification.CreateRouter)
				notificationRouter.PUT("/update/:id", notification.UpdateRouter)
				notificationRouter.DELETE("/delete/:id", notification.DeleteRouter)
			}

			// 个人消息
			messageRouter := adminRouter.Group("/message")
			{
				messageRouter.POST("/create", message.CreateRouter)
				messageRouter.PUT("/update/:message_id", message.UpdateRouter)
				messageRouter.DELETE("/delete/:message_id", message.DeleteByAdminRouter)
			}

			// Banner
			bannerRouter := adminRouter.Group("banner")
			{
				bannerRouter.POST("/create", banner.CreateRouter)
				bannerRouter.PUT("/update/:banner_id", banner.UpdateRouter)
			}
		}

		// 认证类
		authRouter := v1.Group("/auth")
		{
			authRouter.POST("/signup", auth.SignUpRouter)
			authRouter.POST("/signin", auth.SignInRouter)
			authRouter.POST("/activation", auth.ActivationRouter)
			authRouter.PUT("/password/reset", auth.ResetPasswordRouter)

		}

		oauthRouter := v1.Group("/oauth2")
		{
			oauthRouter.GET("/google", oauth2.GoogleLoginRouter)
			oauthRouter.GET("/google_callback", oauth2.GoogleCallbackRouter)
		}

		// 用户类
		userRouter := v1.Group("/user")
		{
			userRouter.Use(userAuthMiddleware)
			userRouter.GET("/signout", user.SignOut)                             // 用户登出
			userRouter.GET("/profile", user.GetProfileRouter)                    // 获取用户详细信息
			userRouter.PUT("/profile", user.UpdateProfileRouter)                 // 更新用户资料
			userRouter.PUT("/password", user.UpdatePasswordRouter)               // 更新登陆密码
			userRouter.POST("/password2", user.SetPayPasswordRouter)             // 设置交易密码
			userRouter.PUT("/password2", user.UpdatePayPasswordRouter)           // 更新交易密码
			userRouter.PUT("/password2/reset", user.ResetPayPasswordRouter)      // 重置交易密码
			userRouter.POST("/password2/reset", user.SendResetPayPasswordRouter) // 发送重置交易密码的邮件/短信
			userRouter.POST("/avatar", user.UploadAvatarRouter)                  // 上传用户头像
			// 邀请人列表
			inviteRouter := userRouter.Group("/invite")
			{
				inviteRouter.GET("/detail/:invite_id", invite.GetRouter) // 获取单条邀请记录详情
				inviteRouter.GET("/list", invite.GetListRouter)          // 获取我已邀请的列表
			}
			// 收货地址
			addressRouter := userRouter.Group("/address")
			{
				addressRouter.POST("/create", address.CreateRouter)               // 添加收货地址
				addressRouter.PUT("/update/:address_id", address.UpdateRouter)    // 更新收货地址
				addressRouter.DELETE("/delete/:address_id", address.DeleteRouter) // 删除收货地址
				addressRouter.GET("/detail/:address_id", address.GetDetailRouter) // 获取地址详情
				addressRouter.GET("/list", address.GetListRouter)                 // 获取地址列表
				addressRouter.GET("/default", address.GetDefaultRouter)           // 获取默认地址
			}
		}

		// 钱包类
		walletRouter := v1.Group("/wallet")
		{
			walletRouter.Use(userAuthMiddleware)
			// 获取所有的钱包信息
			walletRouter.GET("/map", wallet.GetWalletsRouter)               // 获取我的钱包map对象
			walletRouter.GET("/currency/:currency", wallet.GetWalletRouter) // 获取单个钱包的详细信息
			// 转账相关
			walletRouter.GET("/transfer/history", transfer.GetHistory)                    // 获取转账记录
			walletRouter.GET("/transfer/detail/:transfer_id", transfer.GetDetailRouter)   // 获取单条转账详情
			walletRouter.POST("/transfer", middleware.AuthPayPassword, transfer.ToRouter) // 转账给某人
		}

		// 财务日志
		financeRouter := v1.Group("/finance")
		{
			financeRouter.Use(userAuthMiddleware)
			financeRouter.GET("/history", finance.GetHistory) // TODO: 获取我的财务日志
		}

		// 新闻咨询类
		newsRouter := v1.Group("/news")
		{
			newsRouter.GET("/list", news.GetListRouter)
			newsRouter.GET("/detail/:id", news.GetNewsRouter)
		}

		// 系统通知
		notificationRouter := v1.Group("/notification")
		{
			walletRouter.Use(userAuthMiddleware)
			notificationRouter.GET("/list", notification.GetListRouter)
			notificationRouter.GET("/detail/:id", notification.GetRouter)
			notificationRouter.GET("/read/:id", notification.ReadRouter)
		}

		// 用户的个人通知, 用人通知是可以删除的
		messageRouter := v1.Group("/message")
		{
			walletRouter.Use(userAuthMiddleware)
			messageRouter.GET("/list", message.GetListRouter)
			messageRouter.GET("/detail/:id", message.GetRouter)
			messageRouter.GET("/read/:id", message.ReadRouter)
			messageRouter.DELETE("/delete/:id", message.DeleteByUserRouter)
		}

		// 通用类
		{
			// 邮件服务
			v1.POST("/email/send/activation", email.SendActivationEmailRouter)
			v1.POST("/email/send/reset_password", email.SendResetPasswordEmailRouter)

			// 文件上传 (需要验证token)
			uploadRouter := v1.Group("/upload")
			{
				uploadRouter.Use(userAuthMiddleware)
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
			v1.GET("/avatar/:filename", user.GetAvatarRouter) // 获取用户头像

			v1.GET("/area", address.AreaListRouter) // 获取地址选择列表
		}
	}

	Router = router
}
