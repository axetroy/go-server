package worker

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
)

// 处理来之用户端的消息
func MessageFromUserHandler() {
	for {
		// 从客服池中取消息
		msg := <-ws.WaiterPoll.Broadcast

	typeCondition:
		switch ws.TypeResponseWaiter(msg.Type) {
		// 发送数据给客服
		case ws.TypeResponseWaiterMessageText:
			waiterClient := ws.WaiterPoll.Get(msg.To)

			_ = waiterClient.WriteJSON(ws.Message{
				From:    msg.From,
				To:      msg.To,
				Type:    msg.Type,
				Payload: msg.Payload,
			})
			break typeCondition
		default:
			break typeCondition
		}
	}
}
