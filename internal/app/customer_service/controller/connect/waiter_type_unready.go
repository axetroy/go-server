// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"time"
)

func waiterTypeUnReadyHandler(waiterClient *ws.Client) (err error) {
	if waiterClient.GetProfile() == nil {
		err = exception.UserNotLogin
		return
	}

	waiterClient.SetReady(false)

	// 给予回执
	_ = waiterClient.WriteJSON(ws.Message{
		From: waiterClient.UUID,
		To:   waiterClient.UUID,
		Type: string(ws.TypeResponseWaiterUnready),
		Date: time.Now().Format(time.RFC3339Nano),
	})

	return nil
}
