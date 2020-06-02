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

var UserRouter = router.Handler(func(c router.Context) {
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
			waiterId := ws.MatcherPool.GetMyWaiter(client.UUID)

			// 通知客服，我断开连接
			if waiterId != nil {
				waiterClient := ws.WaiterPoll.Get(*waiterId)

				_ = waiterClient.WriteJSON(ws.Message{
					From:    client.UUID,
					To:      *waiterId,
					Type:    string(ws.TypeToWaiterDisconnected),
					Payload: nil,
				})
			}

			// 从池中删除该链接
			ws.UserPoll.Remove(client.UUID)

			// 断开匹配
			ws.MatcherPool.Leave(client.UUID)

			// 因为当前连接已经断开，应该会空出一个位置
			// 让客服继续接待下一个
			ws.MatcherPool.Broadcast <- true
		}
	}()

	client = ws.NewClient(webscoket)

	// 注册新的客户端
	ws.UserPoll.Add(client)

	// 告诉客户端它的 UUID
	_ = client.WriteJSON(ws.Message{
		Type: string(ws.TypeToUserInitialize),
		To:   client.UUID,
	})

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

	typeCondition:
		switch ws.TypeFromUser(msg.Type) {
		// 连接一个客服
		case ws.TypeFromUserConnect:
			waiterID := ws.MatcherPool.Join(client.UUID)

			// 如果找不到合适的客服，则添加到等待队列
			if waiterID == nil {
				// 告诉客户端要排队
				_ = client.WriteJSON(ws.Message{
					Type: string(ws.TypeToUserConnectQueue),
					To:   client.UUID,
				})
				break typeCondition
			}

			// 告诉客户端已连接成功
			_ = client.WriteJSON(ws.Message{
				Type: string(ws.TypeToUserConnectSuccess),
				From: *waiterID,
				To:   client.UUID,
				Payload: map[string]interface{}{
					"uuid": *waiterID,
					// TODO: 服务的基本信息
				},
			})

			// 告诉客服有新的连接接入
			waiterClient := ws.WaiterPoll.Get(*waiterID)

			if waiterClient != nil {
				_ = waiterClient.WriteJSON(ws.Message{
					From: client.UUID,
					To:   *waiterID,
					Type: string(ws.TypeToWaiterNewConnection),
				})
			}

			break typeCondition
		case ws.TypeFromUserDisconnect:
			ws.MatcherPool.Leave(client.UUID)
			waiterId := ws.MatcherPool.GetMyWaiter(client.UUID)

			// 通知客服，我断开连接
			if waiterId != nil {
				waiterClient := ws.WaiterPoll.Get(*waiterId)

				_ = waiterClient.WriteJSON(ws.Message{
					From:    client.UUID,
					To:      *waiterId,
					Type:    string(ws.TypeToWaiterDisconnected),
					Payload: nil,
				})
			}
			break typeCondition
		// 用户发送消息
		case ws.TypeFromUserMessageText:
			waiterId := ws.MatcherPool.GetMyWaiter(client.UUID)

			// 如果这个客户端没有连接客服，那么消息不会发送
			if waiterId != nil {
				// 把收到的消息广播到客服池
				ws.WaiterPoll.Broadcast <- ws.Message{
					From:    client.UUID,
					Type:    msg.Type,
					To:      *waiterId,
					Payload: msg.Payload,
				}
			} else {
				_ = client.WriteJSON(ws.Message{
					To:   client.UUID,
					Type: string(ws.TypeToUserNotConnect),
				})
			}
			break typeCondition
		default:
			break typeCondition
		}
	}
})
