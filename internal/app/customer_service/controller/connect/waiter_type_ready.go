// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
)

func waiterTypeReadyHandler(waiterClient *ws.Client) error {
	if waiterClient.GetProfile() == nil {
		return exception.UserNotLogin
	}

	// 添加客服到池里，并且分配正在排队的用户
	ws.MatcherPool.AddWaiter(waiterClient.UUID)

	// 获取客服要服务的用户
	users := ws.MatcherPool.GetMyUsers(waiterClient.UUID)

	// 把正在排队的用户，分配给这个客服
	for _, userSocketUUID := range users {
		userClient := ws.UserPoll.Get(userSocketUUID)
		if userClient != nil {
			// 告诉用户端已连接成功
			_ = userClient.WriteJSON(ws.Message{
				From:    waiterClient.UUID,
				To:      userSocketUUID,
				Type:    string(ws.TypeResponseUserConnectSuccess),
				Payload: waiterClient.GetProfile(),
			})
			// 告诉客服端已连接成功
			_ = waiterClient.WriteJSON(ws.Message{
				From:    userSocketUUID,
				To:      waiterClient.UUID,
				Type:    string(ws.TypeResponseWaiterNewConnection),
				Payload: userClient.GetProfile(),
			})
		}
	}

	return nil
}
