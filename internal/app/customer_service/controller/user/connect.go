package user

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

var ConnectRouter = router.Handler(func(c router.Context) {
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
			// 从池中删除该链接
			ws.UserPoll.RemoveClient(client.UUID)
			// 断开匹配
			ws.MatcherPool.Disconnect(client.UUID)
			// 从队列中删除
			ws.MatcherPool.RemoveWaiter(client.UUID)
			// TODO: 通知已连接的客服断开连接
		}
	}()

	client = ws.NewClient(webscoket)

	// 注册新的客户端
	ws.UserPoll.AddClient(client)

	// 告诉客户端它的 UUID
	_ = client.WriteJSON(ws.Message{
		Type: string(ws.TypeToUserInitialize),
		To:   &client.UUID,
	})

	for {
		var msg ws.Message
		// 读取消息
		err := webscoket.ReadJSON(&msg)

		if err != nil {
			_ = client.WriteJSON(ws.Message{
				Type: string(ws.TypeToUserError),
				To:   &client.UUID,
				Payload: map[string]string{
					"message": exception.InvalidParams.Error(),
				},
			})
			continue
		}

		// 传入的参数不正确
		if err := validator.ValidateStruct(msg); err != nil {
			_ = client.WriteJSON(ws.Message{
				Type: string(ws.TypeToUserError),
				To:   &client.UUID,
				Payload: map[string]string{
					"message": exception.InvalidParams.Error(),
				},
			})
			continue
		}

	typeCondition:
		switch ws.TypeFromUser(msg.Type) {
		// 连接一个客服
		case ws.TypeFromUserConnect:
			waiterID := ws.MatcherPool.LookupWaiter()

			// 如果找不到合适的客服，则添加到等待队列
			if waiterID == nil {
				ws.MatcherPool.AppendToQueue(client.UUID)
				// 告诉客户端要排队
				_ = client.WriteJSON(ws.Message{
					Type: string(ws.TypeToUserConnectQueue),
					To:   &client.UUID,
				})
				break typeCondition
			}

			ws.MatcherPool.Connect(*waiterID, client.UUID)

			// 告诉客户端已连接成功
			_ = client.WriteJSON(ws.Message{
				Type: string(ws.TypeToUserConnectSuccess),
				From: *waiterID,
				To:   &client.UUID,
				Payload: map[string]string{
					"uuid": *waiterID,
					// TODO: 服务的基本信息
				},
			})

			// 告诉客服有新的连接接入
			waiterClient := ws.WaiterPoll.GetClient(*waiterID)

			if waiterClient != nil {
				_ = waiterClient.WriteJSON(ws.Message{
					To:   waiterID,
					Type: string(ws.TypeToWaiterNewConnection),
					Payload: map[string]string{
						"uuid": client.UUID,
					},
				})
			}

			break typeCondition
		// 用户发送消息
		case ws.TypeFromUserMessageText:
			waiterId := ws.MatcherPool.GetCurrentWaiter(client.UUID)

			// 如果这个客户端没有连接客服，那么消息不会发送
			if waiterId != nil {
				// 把收到的消息广播到客服池
				ws.WaiterPoll.Broadcast <- ws.Message{
					From:    client.UUID,
					Type:    msg.Type,
					To:      waiterId,
					Payload: msg.Payload,
				}
			}
			break typeCondition
		default:
			break typeCondition
		}
	}
})
