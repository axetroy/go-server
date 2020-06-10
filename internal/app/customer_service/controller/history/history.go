// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package history

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"time"
)

type History struct {
	ID       string               `json:"id"`       // 消息 ID
	Sender   schema.ProfilePublic `json:"sender"`   // 消息发送者
	Receiver schema.ProfilePublic `json:"receiver"` // 消息接受者
	Type     ws.TypeResponseUser  `json:"type"`     // 消息类型
	Payload  interface{}          `json:"payload"`  // 消息体
	Date     string               `json:"date"`     // 消息时间
}

// 获取某个用户的聊天记录
func GetHistory(userID string, txs ...*gorm.DB) (result []History, err error) {
	var tx *gorm.DB
	if len(txs) > 0 {
		tx = txs[0]
	}

	if tx == nil {
		tx = database.Db.Begin()
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback().Error
		} else {
			_ = tx.Commit().Error
		}
	}()

	list := make([]model.CustomerSessionItem, 0)

	query := tx.Model(model.CustomerSessionItem{}).Where("sender_id = ?", userID).Or("receiver_id = ?", userID).Order("created_at ASC").Preload("Sender").Preload("Receiver").Limit(100)

	if err = query.Find(&list).Error; err != nil {
		return
	}

	for _, info := range list {
		target := History{
			ID: info.Id,
			Sender: schema.ProfilePublic{
				Id:       info.Sender.Id,
				Username: info.Sender.Username,
				Nickname: info.Sender.Nickname,
				Avatar:   info.Sender.Avatar,
			},
			Receiver: schema.ProfilePublic{
				Id:       info.Receiver.Id,
				Username: info.Receiver.Username,
				Nickname: info.Receiver.Nickname,
				Avatar:   info.Receiver.Avatar,
			},
			Payload: info.Payload,
			Date:    info.CreatedAt.Format(time.RFC3339Nano),
		}

		switch info.Type {
		case model.SessionTypeText:
			target.Type = ws.TypeResponseUserMessageText

			type Payload struct {
				Message string `json:"message"`
			}

			var payload Payload
			if err = json.Unmarshal([]byte(info.Payload), &payload); err != nil {
				return
			}

			target.Payload = payload
		case model.SessionTypeImage:
		}

		result = append(result, target)
	}

	return
}
