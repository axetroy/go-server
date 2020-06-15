// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
)

func waiterTypeReadyHandler(waiterClient *ws.Client) (err error) {
	if waiterClient.GetProfile() == nil {
		err = exception.UserNotLogin
		return
	}

	waiterClient.SetReady(true)

	ws.MatcherPool.AddWaiter(waiterClient.UUID)

	i := 0

	// 让这个客服接满客
	for {
		if i > ws.MatcherPool.Max {
			break
		}
		// 通知匹配池，开始匹配
		ws.MatcherPool.Broadcast <- true
		i = i + 1
	}

	return err
}
