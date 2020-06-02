// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package customer_service

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/internal/app/customer_service/controller/user"
	"github.com/axetroy/go-server/internal/app/customer_service/controller/waiter"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"github.com/kataras/iris/v12/middleware/recover"
	"log"
	"net/http"
)

var (
	CustomerServiceRouter *iris.Application
)

func handleMessagesFromUser() {
	for {
		// 从客服池中取消息
		msg := <-ws.WaiterPoll.Broadcast

		waiterID := ws.MatcherPool.GetCurrentWaiter(msg.From)

		if waiterID == nil {
			continue
		}

	typeCondition:
		switch ws.TypeToWaiter(msg.Type) {
		// 发送数据给客服
		case ws.TypeToWaiterMessageText:
			client := ws.WaiterPoll.GetClient(*waiterID)

			err := client.WriteJSON(ws.Message{
				From:    msg.From,
				To:      waiterID,
				Type:    msg.Type,
				Payload: msg.Payload,
			})
			// TODO: 处理发送失败的情况
			if err != nil {
				log.Printf("error: %v\n", err)
			}
			break typeCondition
		default:
			break typeCondition
		}
	}
}

func handleMessagesFromWaiter() {
	for {
		// 从用户池中取消息
		msg := <-ws.UserPoll.Broadcast

		// 如果没有指明发给谁，那么跳过
		if msg.To == nil {
			continue
		}

		// 如果用户已经不匹配了，那么不处理这条消息
		if ws.MatcherPool.IsUserConnectingWithWaiter(msg.From, *msg.To) == false {
			continue
		}

	typeCondition:
		switch ws.TypeToWaiter(msg.Type) {
		// 发送消息给用户
		case ws.TypeToWaiterMessageText:

			client := ws.UserPoll.GetClient(*msg.To)

			err := client.WriteJSON(ws.Message{
				From:    msg.From,
				To:      msg.To,
				Type:    string(ws.TypeToUserMessageText),
				Payload: msg.Payload,
			})
			// TODO: 处理发送失败的情况
			if err != nil {
				log.Printf("error: %v\n", err)
			}
			break typeCondition
		default:
			break typeCondition
		}
	}
}

func init() {
	app := iris.New()

	app.OnAnyErrorCode(router.Handler(func(c router.Context) {
		code := c.GetStatusCode()

		c.StatusCode(code)

		c.JSON(errors.New(fmt.Sprintf("%d %s", code, http.StatusText(code))), nil, nil)
	}))

	v1 := app.Party("v1")

	{
		v1.Use(recover.New())
		v1.Use(middleware.Common())
		v1.Use(middleware.CORS())

		if config.Common.Mode != "production" {
			v1.Use(logger.New())
			v1.Use(middleware.Ip())
		}

		// 接收监听的消息并进行处理
		go handleMessagesFromUser()
		go handleMessagesFromWaiter()

		{
			// 连接 WebSocket
			socketRouter := v1.Party("/ws")
			socketRouter.Get("/user/connect", user.ConnectRouter)     // 用户连接
			socketRouter.Get("/waiter/connect", waiter.ConnectRouter) // 客服连接
		}
	}

	_ = app.Build()

	CustomerServiceRouter = app
}
