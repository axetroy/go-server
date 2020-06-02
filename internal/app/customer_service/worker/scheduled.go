package worker

import "github.com/axetroy/go-server/internal/app/customer_service/ws"

// 任务分配调度器
// 用于分配空闲的客服和正在排队的用户
func DistributionSchedulerHandler() {
	for {
		// 从客服池中取消息
		_ = <-ws.MatcherPool.Broadcast

		idleWaiter := ws.MatcherPool.GetIdleWaiter()

		if idleWaiter != nil {
			userSocketUUID := ws.MatcherPool.ShiftPending()
			if userSocketUUID != nil {
				ws.MatcherPool.Join(*userSocketUUID)

				// 通知双方连接
				userClient := ws.UserPoll.Get(*userSocketUUID)
				waiterClient := ws.WaiterPoll.Get(*idleWaiter)

				_ = userClient.WriteJSON(ws.Message{
					From: *idleWaiter,
					To:   *userSocketUUID,
					Type: string(ws.TypeToUserConnectSuccess),
				})

				_ = waiterClient.WriteJSON(ws.Message{
					From: *userSocketUUID,
					To:   *idleWaiter,
					Type: string(ws.TypeToUserConnectSuccess),
				})
			}
		}
	}
}
