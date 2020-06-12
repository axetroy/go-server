// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"time"
)

func waiterTypeMessageHandler(waiterClient *ws.Client, msg ws.Message) (err error) {
	type MessageBody struct {
		Message string `json:"message" validate:"required" comment:"消息体"` // 发送的消息体
	}

	if waiterClient.GetProfile() == nil {
		err = exception.UserNotLogin
		return err
	}

	var body MessageBody

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
	// 把收到的消息发送给用户
	ws.UserPoll.Broadcast <- ws.Message{
		From:    waiterClient.UUID,
		To:      msg.To,
		Type:    msg.Type,
		Payload: msg.Payload,
		Date:    time.Now().Format(time.RFC3339Nano),
	}
	return nil
}
