// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"time"
)

func waiterTypeMessageImageHandler(waiterClient *ws.Client, msg ws.Message) (err error) {
	// 如果还没有认证
	if waiterClient.GetProfile() == nil {
		return exception.UserNotLogin
	}

	var body ws.MessageImagePayload

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

	// 如果这个客户端没有连接客服，那么消息不会发送
	ws.UserPoll.Broadcast <- ws.Message{
		Type:    msg.Type,
		From:    waiterClient.UUID,
		To:      msg.To,
		Payload: msg.Payload,
		Date:    time.Now().Format(time.RFC3339Nano),
	}

	return err
}
