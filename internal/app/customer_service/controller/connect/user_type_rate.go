// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package connect

import (
	"errors"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"time"
)

func userTypeRateHandler(userClient *ws.Client, msg ws.Message) (err error) {
	// 如果还没有认证
	if userClient.GetProfile() == nil {
		return exception.UserNotLogin
	}

	var body ws.RatePayload

	if err = util.Decode(&body, msg.Payload); err != nil {
		return err
	}

	if err = validator.ValidateStruct(&body); err != nil {
		return err
	}

	waiterId := ws.MatcherPool.GetMyWaiter(userClient.UUID)

	if waiterId == nil {
		return errors.New("未连接")
	}

	waiterClient := ws.WaiterPoll.Get(*waiterId)

	if waiterClient != nil {
		var tx = database.Db.Begin()

		defer func() {
			if err != nil {
				err = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}()

		sessionID := util.MD5(userClient.UUID + waiterClient.UUID)

		// 更新评分
		if err = tx.Model(model.CustomerSession{}).
			Where("id = ?", sessionID).
			Where("uid = ?", userClient.GetProfile().Id).
			Where("rate = NULL").
			Where("closed_at = NULL").
			Update("rate = ?", body.Rate).
			Error; err != nil {
			return
		}

		now := time.Now().Format(time.RFC3339Nano)

		// 给用户回执
		_ = userClient.WriteJSON(ws.Message{
			To:      userClient.UUID,
			Type:    ws.TypeResponseUserRateSuccess.String(),
			Payload: body,
			Date:    now,
		})

		// 给客服回执
		_ = waiterClient.WriteJSON(ws.Message{
			From:    userClient.UUID,
			To:      waiterClient.UUID,
			Type:    ws.TypeResponseWaiterRateUserSuccess.String(),
			Payload: body,
			Date:    now,
		})

		// 断开 socket 连接
		_ = userClient.Close()
	} else {
		return errors.New("未连接")
	}

	return err
}
