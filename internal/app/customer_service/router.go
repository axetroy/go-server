// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package customer_service

import (
	"fmt"
	"github.com/axetroy/go-server/internal/app/customer_service/controller/connect"
	"github.com/axetroy/go-server/internal/app/customer_service/controller/example"
	"github.com/axetroy/go-server/internal/app/customer_service/controller/status"
	"github.com/axetroy/go-server/internal/app/customer_service/worker"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"net/http"
)

var (
	CustomerServiceRouter *iris.Application
)

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
		v1.Use(middleware.RateLimit(30))

		if config.Common.Mode != "production" {
			v1.Use(logger.New())
			v1.Use(middleware.Ip())
		}

		go worker.MessageFromUserHandler()       // 监听来自用户发来的消息
		go worker.MessageFromWaiterHandler()     // 监听来之客服发来的消息
		go worker.DistributionSchedulerHandler() // 调度器

		{
			// 连接 WebSocket
			wsRouter := v1.Party("/ws")

			wsRouter.Get("/status", status.GetStatusRouter) // 客服连接

			{
				connectRouter := wsRouter.Party("/connect")
				connectRouter.Get("/user", connect.UserRouter)     // 用户连接
				connectRouter.Get("/waiter", connect.WaiterRouter) // 客服连接
			}

			{
				connectRouter := wsRouter.Party("/example")
				connectRouter.Get("/user", example.UserRouter)     // 用户端的示例
				connectRouter.Get("/waiter", example.WaiterRouter) // 服务端的示例
			}

		}
	}

	_ = app.Build()

	CustomerServiceRouter = app
}
