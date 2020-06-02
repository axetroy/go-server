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
						Type:    string(ws.TypeToUserDisconnected),
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
	ws.MatcherPool.AddWaiter(client.UUID)

	// 告诉客户端它的 UUID
	_ = client.WriteJSON(ws.Message{
		Type: string(ws.TypeToWaiterTypeInitializeToUser),
		To:   client.UUID,
	})

	users := ws.MatcherPool.GetMyUsers(client.UUID)

	for _, userSocketUUID := range users {
		userClient := ws.UserPoll.Get(userSocketUUID)
		if userClient != nil {
			// 告诉用户端已连接成功
			_ = userClient.WriteJSON(ws.Message{
				From: client.UUID,
				To:   userSocketUUID,
				Type: string(ws.TypeToUserConnectSuccess),
			})
			// 告诉客服端已连接成功
			_ = client.WriteJSON(ws.Message{
				From: userSocketUUID,
				To:   client.UUID,
				Type: string(ws.TypeToWaiterNewConnection),
			})
		}

	}

	for {
		var msg ws.Message
		// 读取消息
		err := webscoket.ReadJSON(&msg)

		if err != nil {
			_ = client.WriteJSON(ws.Message{
				Type: string(ws.TypeToUserError),
				To:   client.UUID,
				Payload: map[string]interface{}{
					"message": exception.InvalidParams.New(err.Error()).Error(),
					"status":  exception.InvalidParams.Code(),
					"data":    msg,
				},
			})
			continue
		}

		// 传入的参数不正确
		if err := validator.ValidateStruct(msg); err != nil {
			_ = client.WriteJSON(ws.Message{
				Type: string(ws.TypeToWaiterError),
				To:   client.UUID,
				Payload: map[string]interface{}{
					"message": exception.InvalidParams.New(err.Error()).Error(),
					"status":  exception.InvalidParams.Code(),
					"data":    msg,
				},
			})
		}

	typeCondition:
		switch ws.TypeFromWaiter(msg.Type) {
		case ws.TypeFromWaiterReady:
			break typeCondition
		case ws.TypeFromWaiterMessageText:
			// 如果没有指定发送给谁
			if msg.To == "" {
				_ = client.WriteJSON(ws.Message{
					Type: string(ws.TypeToUserError),
					To:   client.UUID,
					Payload: map[string]interface{}{
						"message": exception.InvalidParams.Error(),
						"status":  exception.InvalidParams.Code(),
						"data":    msg,
					},
				})
				break typeCondition
			}
			// 把收到的消息发送给用户
			ws.UserPoll.Broadcast <- ws.Message{
				From:    client.UUID,
				To:      msg.To,
				Type:    msg.Type,
				Payload: msg.Payload,
			}
			break typeCondition
		default:
			break typeCondition
		}
	}
})
