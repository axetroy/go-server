// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/token"
	"time"
)

func waiterTypeAuthHandler(waiterClient *ws.Client, msg ws.Message) (err error) {
	type AuthBody struct {
		Token string `json:"token" validate:"required" comment:"Token"`
	}

	var body AuthBody

	if err = util.Decode(&body, msg.Payload); err != nil {
		return err
	}

	if err = validator.ValidateStruct(&body); err != nil {
		return err
	}

	c, err := token.Parse(body.Token, token.StateUser)

	if err != nil {
		return err
	}

	userInfo := model.User{
		Id: c.Uid,
	}

	if err = database.Db.Model(&userInfo).Where(&userInfo).Where("role @> ARRAY[?::varchar]", "waiter").First(&userInfo).Error; err != nil {
		return err
	}

	var profile schema.ProfilePublic

	if err = util.Decode(&profile, userInfo); err != nil {
		return err
	}

	waiterClient.UpdateProfile(profile)

	// 如果这个客服之前已经登录，那么我们就把原有的连接关闭
	if waiters := ws.WaiterPoll.GetWaiterFromUserID(profile.Id); len(waiters) > 0 {
		for _, c := range waiters {
			// 其他连接都要关闭
			if c.UUID != waiterClient.UUID {
				// 推送断开连接
				if err = c.WriteJSON(ws.Message{
					Type: string(ws.TypeResponseWaiterKickOut),
					From: waiterClient.UUID,
					To:   c.UUID,
					Date: time.Now().Format(time.RFC3339Nano),
				}); err != nil {
					return
				}
				_ = c.Close()
			}
		}
	}

	// 告诉客服端它的身份信息
	if err = waiterClient.WriteJSON(ws.Message{
		Type:    string(ws.TypeResponseWaiterAuthSuccess),
		To:      waiterClient.UUID,
		Payload: profile,
		Date:    time.Now().Format(time.RFC3339Nano),
	}); err != nil {
		return
	}

	return err
}
