// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/gorilla/websocket"
	"net/http"
)

var WaiterRouter = router.Handler(func(c router.Context) {
	var (
		client *ws.Client
	)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	webscoket, err := upgrader.Upgrade(c.Writer(), c.Request(), nil)

	if err != nil {
		c.ResponseFunc(nil, func() schema.Response {
			return schema.Response{
				Message: http.StatusText(http.StatusInternalServerError),
				Status:  schema.StatusFail,
				Data:    nil,
			}
		})
		return
	}

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}

			fmt.Printf("%+v\n", err)
		}

		_ = webscoket.Close()
		if client != nil {
			// 移除客户端
			ws.WaiterPoll.Remove(client.UUID)
			// 通知已连接的用户断开连接
			users := ws.MatcherPool.GetMyUsers(client.UUID)

			for _, user := range users {
				userSocket := ws.UserPoll.Get(user)
				if userSocket != nil {
					_ = userSocket.WriteJSON(ws.Message{
						From:    client.UUID,
						To:      user,
						Type:    string(ws.TypeResponseUserDisconnected),
						Payload: nil,
					})
				}
			}
			// 从客服池中移除
			ws.MatcherPool.RemoveWaiter(client.UUID)

			// 因为当前连接已经断开，正在连接的用户会被加入到队列
			// 所以触发一次任务调度
			ws.MatcherPool.Broadcast <- true
		}
	}()

	client = ws.NewClient(webscoket)

	// 注册新的客户端
	ws.WaiterPoll.Add(client)

	for {
		var msg ws.Message
		// 读取消息
		err := webscoket.ReadJSON(&msg)

		if err != nil {
			_ = client.WriteError(exception.InvalidParams.New(err.Error()), msg)
			continue
		}

		// 传入的参数不正确
		if err := validator.ValidateStruct(msg); err != nil {
			_ = client.WriteError(exception.InvalidParams.New(err.Error()), msg)
			continue
		}

	typeCondition:
		switch ws.TypeRequestWaiter(msg.Type) {
		case ws.TypeRequestWaiterAuth:
			if er := waiterTypeAuthHandler(client, msg); er != nil {
				_ = client.WriteError(er, msg)
			}
			break typeCondition
		case ws.TypeRequestWaiterReady:
			if er := waiterTypeReadyHandler(client); er != nil {
				_ = client.WriteError(er, msg)
			}
			break typeCondition
		case ws.TypeRequestWaiterDisconnect:
			if er := waiterTypeDisconnectHandler(client, msg); er != nil {
				_ = client.WriteError(er, msg)
			}
			break typeCondition
		case ws.TypeRequestWaiterMessageText:
			if er := waiterTypeMessageHandler(client, msg); er != nil {
				_ = client.WriteError(exception.InvalidParams.New(er.Error()), msg)
			}
			break typeCondition
		default:
			_ = client.WriteError(exception.InvalidParams.New("未知的消息类型"), msg)
			break typeCondition
		}
	}
})
