// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package resource_server

import (
	"fmt"
	"github.com/axetroy/go-server/internal/app/resource_server/controller/downloader"
	"github.com/axetroy/go-server/internal/app/resource_server/controller/resource"
	"github.com/axetroy/go-server/internal/app/resource_server/controller/uploader"
	"github.com/axetroy/go-server/internal/app/user_server/controller/user"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/dotenv"
	"github.com/gin-gonic/gin"
	"net/http"
	"path"
)

var ResourceRouter *gin.Engine

func init() {
	if config.Common.Mode == config.ModeProduction {
		gin.SetMode(gin.ReleaseMode)
	}
	router := gin.Default()

	router.Use(middleware.GracefulExit())

	router.Use(middleware.CORS())

	router.Static("/public", path.Join(dotenv.RootDir, "public"))

	if config.Common.Mode != config.ModeProduction {
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

		// 通用类
		{
			// 文件上传
			v1.POST("/upload/file", uploader.File)      // 上传文件
			v1.POST("/upload/image", uploader.Image)    // 上传图片
			v1.GET("/upload/example", uploader.Example) // 上传文件的 example
			// 单纯获取资源文本
			v1.GET("/resource/file/:filename", resource.File)           // 获取文件纯文本
			v1.GET("/resource/image/:filename", resource.Image)         // 获取图片纯文本
			v1.GET("/resource/thumbnail/:filename", resource.Thumbnail) // 获取缩略图纯文本
			// 下载资源
			v1.GET("/download/file/:filename", downloader.File)           // 下载文件
			v1.GET("/download/image/:filename", downloader.Image)         // 下载图片
			v1.GET("/download/thumbnail/:filename", downloader.Thumbnail) // 下载缩略图
			// 公共资源目录
			v1.GET("/avatar/:filename", user.GetAvatarRouter) // 获取用户头像
		}

	}

	ResourceRouter = router
}
