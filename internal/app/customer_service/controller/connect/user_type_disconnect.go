// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"time"
)

func userTypeDisconnectHandler(userClient *ws.Client) error {
	// 如果还没有认证
	if userClient.GetProfile() == nil {
		return exception.UserNotLogin
	}

	waiterId := ws.MatcherPool.GetMyWaiter(userClient.UUID)

	var fromId string

	// 通知客服，我断开连接
	if waiterId != nil {
		waiterClient := ws.WaiterPoll.Get(*waiterId)

		_ = waiterClient.WriteJSON(ws.Message{
			From:    userClient.UUID,
			To:      *waiterId,
			Type:    string(ws.TypeResponseWaiterDisconnected),
			Payload: nil,
			Date:    time.Now().Format(time.RFC3339Nano),
		})

		fromId = *waiterId
	}

	ws.MatcherPool.Leave(userClient.UUID)

	// 通知自己，连接已断开
	_ = userClient.WriteJSON(ws.Message{
		From:    fromId,
		To:      userClient.UUID,
		Type:    string(ws.TypeResponseUserDisconnected),
		Payload: nil,
		Date:    time.Now().Format(time.RFC3339Nano),
	})

	return nil
}
