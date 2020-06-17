// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/authentication"
	"github.com/axetroy/go-server/internal/service/database"
	"time"
)

func userTypeAuthHandler(userClient *ws.Client, msg ws.Message) (err error) {
	var body ws.AuthPayload

	if err = util.Decode(&body, msg.Payload); err != nil {
		return err
	}

	if err = validator.ValidateStruct(&body); err != nil {
		return err
	}

	uid, err := authentication.Gateway(false).Parse(body.Token)

	if err != nil {
		return
	}

	userInfo := model.User{
		Id: uid,
	}

	if err = database.Db.Model(&userInfo).Where(&userInfo).First(&userInfo).Error; err != nil {
		return err
	}

	var profile schema.ProfilePublic

	if err = util.Decode(&profile, userInfo); err != nil {
		return err
	}

	userClient.UpdateProfile(profile)

	// 告诉客户端它的身份信息
	if err = userClient.WriteJSON(ws.Message{
		Type:    string(ws.TypeResponseUserAuthSuccess),
		To:      userClient.UUID,
		Payload: profile,
		Date:    time.Now().Format(time.RFC3339Nano),
	}); err != nil {
		return
	}

	return err
}
