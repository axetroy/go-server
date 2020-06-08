// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"time"
)

func waiterTypeDisconnectHandler(waiterClient *ws.Client, msg ws.Message) error {
	type DisconnectBody struct {
		UUID string `json:"uuid" validate:"required"`
	}

	var body DisconnectBody

	if err := util.Decode(&body, msg.Payload); err != nil {
		return err
	}

	if err := validator.ValidateStruct(&body); err != nil {
		return err
	}

	userClient := ws.UserPoll.Get(body.UUID)

	if userClient != nil {
		ws.MatcherPool.Leave(body.UUID)

		// 告诉用户端断开连接
		_ = userClient.WriteJSON(ws.Message{
			From:    waiterClient.UUID,
			To:      userClient.UUID,
			Type:    string(ws.TypeResponseUserDisconnected),
			Payload: nil,
			Date:    time.Now().Format(time.RFC3339Nano),
		})

		// 告诉客服端已断开连接
		_ = waiterClient.WriteJSON(ws.Message{
			From:    userClient.UUID,
			To:      waiterClient.UUID,
			Type:    string(ws.TypeResponseWaiterDisconnected),
			Payload: nil,
			Date:    time.Now().Format(time.RFC3339Nano),
		})
	} else {
		return exception.InvalidParams.New("未连接")
	}

	return nil
}
