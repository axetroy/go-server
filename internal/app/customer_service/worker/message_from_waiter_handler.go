package worker

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"log"
)

// 处理来之用户端的消息
func MessageFromWaiterHandler() {
	for {
		// 从用户池中取消息
		msg := <-ws.UserPoll.Broadcast

	typeCondition:
		switch ws.TypeToWaiter(msg.Type) {
		// 发送消息给用户
		case ws.TypeToWaiterMessageText:
			userClient := ws.UserPoll.Get(msg.To)

			err := userClient.WriteJSON(ws.Message{
				From:    msg.From,
				To:      msg.To,
				Type:    string(ws.TypeToUserMessageText),
				Payload: msg.Payload,
			})
			// TODO: 处理发送失败的情况
			if err != nil {
				log.Printf("error: %v\n", err)
			}
			break typeCondition
		default:
			break typeCondition
		}
	}
}
