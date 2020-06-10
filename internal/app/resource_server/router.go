// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package resource_server

import (
	"fmt"
	"github.com/axetroy/go-server/internal/app/resource_server/controller/downloader"
	"github.com/axetroy/go-server/internal/app/resource_server/controller/resource"
	"github.com/axetroy/go-server/internal/app/resource_server/controller/uploader"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"net/http"
)

var ResourceRouter *iris.Application

func init() {
	app := iris.New()

	app.OnAnyErrorCode(router.Handler(func(c router.Context) {
		code := c.GetStatusCode()

		c.StatusCode(code)

		c.JSON(fmt.Errorf("%d %s", code, http.StatusText(code)), nil, nil)
	}))

	v1 := app.Party("v1").AllowMethods(iris.MethodOptions)

	{
		v1.Use(recover.New())
		v1.Use(middleware.Common())
		v1.Use(middleware.CORS())

		if config.Common.Mode != "production" {
			v1.Use(logger.New())
			v1.Use(middleware.Ip())
		}

		{
			v1.Get("", router.Handler(func(c router.Context) {
				c.JSON(nil, map[string]string{"ping": "tong"}, nil)
			}))
		}

		// 通用类
		{
			// 文件上传
			v1.Post("/upload/file", uploader.File)      // 上传文件
			v1.Post("/upload/image", uploader.Image)    // 上传图片
			v1.Get("/upload/example", uploader.Example) // 上传文件的 example
			//// 单纯获取资源文本
			v1.Get("/resource/file/{filename}", resource.File)   // 获取文件纯文本
			v1.Get("/resource/image/{filename}", resource.Image) // 获取图片纯文本
			//// 下载资源
			v1.Get("/download/file/{filename}", downloader.File)   // 下载文件
			v1.Get("/download/image/{filename}", downloader.Image) // 下载图片
		}

	}

	_ = app.Build()

	ResourceRouter = app
}
