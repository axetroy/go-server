// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
)

func waiterTypeMessageHandler(waiterClient *ws.Client, msg ws.Message) error {
	type MessageBody struct {
		Message string `json:"message" valid:"required~请输入消息"`
	}

	if waiterClient.GetProfile() == nil {
		return exception.UserNotLogin
	}

	var body MessageBody

	if err := util.Decode(&body, msg.Payload); err != nil {
		return err
	}

	if err := validator.ValidateStruct(&body); err != nil {
		return err
	}

	// 如果没有指定发送给谁
	if msg.To == "" {
		return exception.InvalidParams.New("缺少发送者")
	}
	// 把收到的消息发送给用户
	ws.UserPoll.Broadcast <- ws.Message{
		From:    waiterClient.UUID,
		To:      msg.To,
		Type:    msg.Type,
		Payload: msg.Payload,
	}
	return nil
}
