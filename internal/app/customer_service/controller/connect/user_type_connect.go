// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"time"

	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
)

func userTypeConnectHandler(userClient *ws.Client, msg ws.Message) (err error) {
	// 如果还没有认证
	if userClient.GetProfile() == nil {
		err = exception.UserNotLogin
		return
	}
	waiterID, location := ws.MatcherPool.Join(userClient.UUID)

	// 如果找不到合适的客服，则添加到等待队列
	if waiterID == nil {
		// 告诉客户端要排队
		if err = userClient.WriteJSON(ws.Message{
			Type: string(ws.TypeResponseUserConnectQueue),
			To:   userClient.UUID,
			Date: time.Now().Format(time.RFC3339Nano),
			Payload: map[string]interface{}{
				"location": location,
			},
		}); err != nil {
			return
		}
		return
	}

	waiterClient := ws.WaiterPoll.Get(*waiterID)

	if waiterClient != nil {

		// 连接成功，那么数据库创建一个会话
		tx := database.Db.Begin()

		defer func() {
			if err != nil {
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}()

		hash := util.MD5(userClient.UUID + waiterClient.UUID)

		session := model.CustomerSession{
			Id:       hash,
			Uid:      userClient.GetProfile().Id,
			WaiterID: waiterClient.GetProfile().Id,
		}

		if err = tx.Create(&session).Error; err != nil {
			return
		}

		// 告诉用户端已连接成功
		if err = userClient.WriteJSON(ws.Message{
			Type:    string(ws.TypeResponseUserConnectSuccess),
			From:    *waiterID,
			To:      userClient.UUID,
			Payload: waiterClient.GetProfile(),
			Date:    time.Now().Format(time.RFC3339Nano),
		}); err != nil {
			return
		}

		// 告诉客服端有新的连接接入
		if err = waiterClient.WriteJSON(ws.Message{
			From:    userClient.UUID,
			To:      *waiterID,
			Type:    string(ws.TypeResponseWaiterNewConnection),
			Payload: userClient.GetProfile(),
			Date:    time.Now().Format(time.RFC3339Nano),
		}); err != nil {
			return
		}
	}

	return
}
