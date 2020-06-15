// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/controller/history"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"time"
)

func waiterTypeGetHistoryHandler(waiterClient *ws.Client, msg ws.Message) (err error) {
	type GetHistoryPayload struct {
		UserID string `json:"user_id" validate:"required" comment:"用户 ID"`
	}

	var body GetHistoryPayload

	if err = util.Decode(&body, msg.Payload); err != nil {
		return err
	}

	if err = validator.ValidateStruct(&body); err != nil {
		return err
	}

	// 如果还没有认证
	if waiterClient.GetProfile() == nil {
		return exception.UserNotLogin
	}

	if list, err := history.GetHistory(body.UserID); err != nil {
		return err
	} else {
		if err = waiterClient.WriteJSON(ws.Message{
			Type: string(ws.TypeResponseWaiterMessageHistory),
			To:   waiterClient.UUID,
			Payload: map[string]interface{}{
				"user_id": body.UserID,
				"data":    list,
			},
			Date: time.Now().Format(time.RFC3339Nano),
		}); err != nil {
			return err
		}
	}

	return err
}

func waiterTypeGetHistorySessionHandler(waiterClient *ws.Client, msg ws.Message) (err error) {
	profile := waiterClient.GetProfile()

	// 如果还没有认证
	if profile == nil {
		return exception.UserNotLogin
	}

	if list, err := history.GetWaiterSession(profile.Id); err != nil {
		return err
	} else {
		if err = waiterClient.WriteJSON(ws.Message{
			Type: string(ws.TypeResponseWaiterSessionHistory),
			To:   waiterClient.UUID,
			Payload: map[string]interface{}{
				"user_id": profile.Id,
				"data":    list,
			},
			Date: time.Now().Format(time.RFC3339Nano),
		}); err != nil {
			return err
		}
	}

	return err
}
