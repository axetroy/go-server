// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/controller/history"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"time"
)

func waiterTypeReadyHandler(waiterClient *ws.Client) (err error) {
	if waiterClient.GetProfile() == nil {
		err = exception.UserNotLogin
		return
	}

	// 添加客服到池里，并且分配正在排队的用户
	ws.MatcherPool.AddWaiter(waiterClient.UUID)

	// 获取客服要服务的用户
	users := ws.MatcherPool.GetMyUsers(waiterClient.UUID)

	// 连接成功，那么数据库创建一个会话
	tx := database.Db.Begin()

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	// 把正在排队的用户，分配给这个客服
	for _, userSocketUUID := range users {
		userClient := ws.UserPoll.Get(userSocketUUID)
		if userClient != nil {
			// 告诉用户端已连接成功
			if err = userClient.WriteJSON(ws.Message{
				From:    waiterClient.UUID,
				To:      userSocketUUID,
				Type:    string(ws.TypeResponseUserConnectSuccess),
				Payload: waiterClient.GetProfile(),
				Date:    time.Now().Format(time.RFC3339Nano),
			}); err != nil {
				return
			}
			// 告诉客服端已连接成功
			if err = waiterClient.WriteJSON(ws.Message{
				From:    userSocketUUID,
				To:      waiterClient.UUID,
				Type:    string(ws.TypeResponseWaiterNewConnection),
				Payload: userClient.GetProfile(),
				Date:    time.Now().Format(time.RFC3339Nano),
			}); err != nil {
				return
			}

			hash := util.MD5(userClient.UUID + waiterClient.UUID)

			session := model.CustomerSession{
				Id:       hash,
				Uid:      userClient.GetProfile().Id,
				WaiterID: waiterClient.GetProfile().Id,
			}

			// 创建 session
			if err = tx.Create(&session).Error; err != nil {
				return
			}

			// 推送聊天记录
			if historyMessage, e := history.GetHistory(userClient.GetProfile().Id); e != nil {
				err = e
				return
			} else {
				// 推送给用户端
				if err = userClient.WriteJSON(ws.Message{
					Type:    string(ws.TypeResponseUserMessageHistory),
					From:    waiterClient.UUID,
					To:      userClient.UUID,
					Payload: historyMessage,
					Date:    time.Now().Format(time.RFC3339Nano),
				}); err != nil {
					return
				}

				// 推送给客服端
				if err = waiterClient.WriteJSON(ws.Message{
					Type:    string(ws.TypeResponseWaiterMessageHistory),
					From:    userClient.UUID,
					To:      waiterClient.UUID,
					Payload: historyMessage,
					Date:    time.Now().Format(time.RFC3339Nano),
				}); err != nil {
					return
				}
			}
		}
	}

	return err
}
