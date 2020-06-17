// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"errors"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"time"
)

func waiterTypeRateHandler(waiterClient *ws.Client, msg ws.Message) (err error) {
	// 如果还没有认证
	if waiterClient.GetProfile() == nil {
		return exception.UserNotLogin
	}

	var body ws.RatePayload

	if err = util.Decode(&body, msg.Payload); err != nil {
		return err
	}

	if err = validator.ValidateStruct(&body); err != nil {
		return err
	}

	// 如果没有指定发送给谁
	if msg.To == "" {
		err = exception.InvalidParams.New("缺少发送者")
		return
	}

	userClient := ws.UserPoll.Get(msg.To)

	if userClient != nil {
		now := time.Now().Format(time.RFC3339Nano)
		// 发给客户
		if err = userClient.WriteJSON(ws.Message{
			Type: ws.TypeResponseUserRate.String(),
			From: waiterClient.UUID,
			To:   userClient.UUID,
			Date: now,
		}); err != nil {
			return
		}

		// 给客服回执
		_ = waiterClient.WriteJSON(ws.Message{
			From:    userClient.UUID,
			To:      waiterClient.UUID,
			Type:    ws.TypeResponseWaiterRateSuccess.String(),
			Payload: body,
			Date:    now,
		})
	} else {
		return errors.New("未连接")
	}

	return err
}
