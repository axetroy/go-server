package worker

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/service/database"
	"log"
	"time"
)

func textMessageFromUserHandler(msg ws.Message) (err error) {
	waiterClient := ws.WaiterPoll.Get(msg.To)
	userClient := ws.UserPoll.Get(msg.From)

	if waiterClient == nil || userClient == nil {
		return
	}

	// 发送成功，写入聊天记录
	tx := database.Db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	hash := util.MD5(userClient.UUID + waiterClient.UUID)

	session := model.CustomerSession{
		Id:       hash,
		Uid:      userClient.GetProfile().Id,
		WaiterID: waiterClient.GetProfile().Id,
	}

	// 获取会话
	if err := tx.Model(model.CustomerSession{}).Where(&session).First(&session).Error; err != nil {
		return err
	}

	var raw []byte

	if raw, err = json.Marshal(msg.Payload); err != nil {
		return
	}

	sessionItem := model.CustomerSessionItem{
		SessionID:  session.Id,
		Type:       model.SessionTypeText,
		ReceiverID: waiterClient.GetProfile().Id,
		SenderID:   userClient.GetProfile().Id,
		Payload:    string(raw),
	}

	// 讲聊天记录写入会话
	if err := tx.Create(&sessionItem).Error; err != nil {
		return err
	}

	// 推送消息给客服
	_ = waiterClient.WriteJSON(ws.Message{
		Id:      session.Id,
		From:    msg.From,
		To:      msg.To,
		Type:    string(ws.TypeResponseWaiterMessageText),
		Payload: msg.Payload,
		Date:    sessionItem.CreatedAt.Format(time.RFC3339Nano),
	})

	// 给用户端一个回执
	_ = userClient.WriteJSON(ws.Message{
		Id:      sessionItem.Id,
		Type:    string(ws.TypeResponseUserMessageTextSuccess),
		From:    userClient.UUID,
		To:      waiterClient.UUID,
		Payload: msg.Payload,
		Date:    sessionItem.CreatedAt.Format(time.RFC3339Nano),
	})

	return
}

func imageMessageFromUserHandler(msg ws.Message) (err error) {
	waiterClient := ws.WaiterPoll.Get(msg.To)
	userClient := ws.UserPoll.Get(msg.From)

	if waiterClient == nil || userClient == nil {
		return
	}

	// 发送成功，写入聊天记录
	tx := database.Db.Begin()

	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	hash := util.MD5(userClient.UUID + waiterClient.UUID)

	session := model.CustomerSession{
		Id:       hash,
		Uid:      userClient.GetProfile().Id,
		WaiterID: waiterClient.GetProfile().Id,
	}

	// 获取会话
	if err := tx.Model(model.CustomerSession{}).Where(&session).First(&session).Error; err != nil {
		return err
	}

	var raw []byte

	if raw, err = json.Marshal(msg.Payload); err != nil {
		return
	}

	sessionItem := model.CustomerSessionItem{
		SessionID:  session.Id,
		Type:       model.SessionTypeImage,
		ReceiverID: waiterClient.GetProfile().Id,
		SenderID:   userClient.GetProfile().Id,
		Payload:    string(raw),
	}

	// 讲聊天记录写入会话
	if err := tx.Create(&sessionItem).Error; err != nil {
		return err
	}

	// 推送消息给客服
	_ = waiterClient.WriteJSON(ws.Message{
		Id:      session.Id,
		From:    msg.From,
		To:      msg.To,
		Type:    string(ws.TypeResponseWaiterMessageImage),
		Payload: msg.Payload,
		Date:    sessionItem.CreatedAt.Format(time.RFC3339Nano),
	})

	// 给用户端一个回执
	_ = userClient.WriteJSON(ws.Message{
		Id:      sessionItem.Id,
		Type:    string(ws.TypeResponseUserMessageImageSuccess),
		From:    userClient.UUID,
		To:      waiterClient.UUID,
		Payload: msg.Payload,
		Date:    sessionItem.CreatedAt.Format(time.RFC3339Nano),
	})

	return
}

// 处理来之用户端的消息
func MessageFromUserHandler() {
	for {
		// 从客服池中取消息
		msg := <-ws.WaiterPoll.Broadcast

	typeCondition:
		switch ws.TypeRequestUser(msg.Type) {
		// 发送文本给客服
		case ws.TypeRequestUserMessageText:
			if err := textMessageFromUserHandler(msg); err != nil {
				log.Println(err)
			}
			break typeCondition
		// 发送图片给客服
		case ws.TypeRequestUserMessageImage:
			if err := imageMessageFromUserHandler(msg); err != nil {
				log.Println(err)
			}
			break typeCondition
		default:
			break typeCondition
		}
	}
}
