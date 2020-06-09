// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package worker

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"log"
	"time"
)

func handle() (err error) {
	idleWaiter := ws.MatcherPool.GetIdleWaiter()
	userSocketUUID := ws.MatcherPool.ShiftPending()

	// 没有排队的用户，那么啥也不干
	if userSocketUUID == nil {
		return nil
	}

	// 虽然有用户还在排队，但是客服不空闲，通知下排队的用户排队的情况
	if idleWaiter == nil {
		// 把用户重新加入到队列中
		ws.MatcherPool.Join(*userSocketUUID, true)
		for location, userID := range ws.MatcherPool.GetPendingQueue() {
			userClient := ws.UserPoll.Get(userID)

			if userClient == nil {
				continue
			}

			// 告诉客户端要排队
			if err = userClient.WriteJSON(ws.Message{
				Type: string(ws.TypeResponseUserConnectQueue),
				To:   userClient.UUID,
				Date: time.Now().Format(time.RFC3339Nano),
				Payload: map[string]interface{}{
					"location": location,
				},
			}); err != nil {
				return
			}
		}
	}

	// 如果有空闲的客服和正在排队的用户，那么就匹配他们
	{
		waiterID, _ := ws.MatcherPool.Join(*userSocketUUID)

		if waiterID == nil {
			return exception.NoData.New("找不到 waiter")
		}

		userClient := ws.UserPoll.Get(*userSocketUUID)
		waiterClient := ws.WaiterPoll.Get(*waiterID)

		if userClient == nil || waiterClient == nil {
			return exception.NoData.New("找不到socket连接")
		}

		// 连接成功，那么数据库创建一个会话
		tx := database.Db.Begin()

		defer func() {
			if err != nil {
				_ = tx.Rollback()
			} else {
				_ = tx.Commit()
			}
		}()

		hash := util.MD5(userClient.UUID + waiterClient.UUID)

		session := model.CustomerSession{
			Id:       hash,
			Uid:      userClient.GetProfile().Id,
			WaiterID: waiterClient.GetProfile().Id,
		}

		// 创建 session
		if err = tx.Create(&session).Error; err != nil {
			return
		}

		if err = userClient.WriteJSON(ws.Message{
			From:    *waiterID,
			To:      *userSocketUUID,
			Type:    string(ws.TypeResponseUserConnectSuccess),
			Payload: userClient.GetProfile(),
		}); err != nil {
			return
		}

		if err = waiterClient.WriteJSON(ws.Message{
			From:    *userSocketUUID,
			To:      *waiterID,
			Type:    string(ws.TypeResponseWaiterNewConnection),
			Payload: userClient.GetProfile(),
		}); err != nil {
			return
		}
	}

	return nil
}

// 任务分配调度器
// 用于分配空闲的客服和正在排队的用户
func DistributionSchedulerHandler() {
	for {
		// 从客服池中取消息
		<-ws.MatcherPool.Broadcast

		if err := handle(); err != nil {
			log.Println(err)
		}
	}
}
