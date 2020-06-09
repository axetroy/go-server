// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/controller/history"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"time"
)

func userTypeGetHistoryHandler(userClient *ws.Client) (err error) {
	// 如果还没有认证
	if userClient.GetProfile() == nil {
		return exception.UserNotLogin
	}

	if list, err := history.GetHistory(userClient.GetProfile().Id); err != nil {
		return err
	} else {
		if err = userClient.WriteJSON(ws.Message{
			Type:    string(ws.TypeResponseUserMessageHistory),
			To:      userClient.UUID,
			Payload: list,
			Date:    time.Now().Format(time.RFC3339Nano),
		}); err != nil {
			return err
		}
	}

	return err
}
