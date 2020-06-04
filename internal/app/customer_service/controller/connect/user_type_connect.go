// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
)

func userTypeConnectHandler(userClient *ws.Client, msg ws.Message) error {
	// 如果还没有认证
	if userClient.GetProfile() == nil {
		return exception.UserNotLogin
	}
	waiterID := ws.MatcherPool.Join(userClient.UUID)

	// 如果找不到合适的客服，则添加到等待队列
	if waiterID == nil {
		// 告诉客户端要排队
		_ = userClient.WriteJSON(ws.Message{
			Type: string(ws.TypeResponseUserConnectQueue),
			To:   userClient.UUID,
		})
		return nil
	}

	waiterClient := ws.WaiterPoll.Get(*waiterID)

	if waiterClient != nil {
		// 告诉用户端已连接成功
		_ = userClient.WriteJSON(ws.Message{
			Type:    string(ws.TypeResponseUserConnectSuccess),
			From:    *waiterID,
			To:      userClient.UUID,
			Payload: waiterClient.GetProfile(),
		})

		// 告诉客服端有新的连接接入
		_ = waiterClient.WriteJSON(ws.Message{
			From:    userClient.UUID,
			To:      *waiterID,
			Type:    string(ws.TypeResponseWaiterNewConnection),
			Payload: userClient.GetProfile(),
		})
	}

	return nil
}
