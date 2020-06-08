// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
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
					Type:    string(ws.TypeResponseWaiterDisconnected),
					Payload: nil,
				})

				hash := util.MD5(client.UUID + waiterClient.UUID)

				now := time.Now()

				// 标记会话为已关闭
				_ = database.Db.Model(model.CustomerSession{}).Where("id = ?", hash).Update(model.CustomerSession{
					ClosedAt: &now,
				}).Error
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
		switch ws.TypeRequestUser(msg.Type) {
		// 身份认证
		case ws.TypeRequestUserAuth:
			if err := userTypeAuthHandler(client, msg); err != nil {
				_ = client.WriteError(err, msg)
			}

			break typeCondition
		// 申请连接一个客服
		case ws.TypeRequestUserConnect:
			if err := userTypeConnectHandler(client, msg); err != nil {
				_ = client.WriteError(err, msg)
			}
			break typeCondition
		// 用户主动和客服断开连接
		case ws.TypeRequestUserDisconnect:
			if err := userTypeDisconnectHandler(client); err != nil {
				_ = client.WriteError(err, msg)
			}
			break typeCondition
		// 用户发送消息
		case ws.TypeRequestUserMessageText:
			if err := userTypeMessageHandler(client, msg); err != nil {
				_ = client.WriteError(err, msg)
			}
			break typeCondition
		default:
			_ = client.WriteError(exception.InvalidParams.New("未知的消息类型"), msg)
			break typeCondition
		}
	}
})
