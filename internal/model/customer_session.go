// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.

package model

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/jinzhu/gorm"
	"time"
)

type SessionType int

const (
	SessionTypeText  SessionType = 0 // 发送纯文本消息
	SessionTypeImage SessionType = 1 // 发送图片
)

var (
	SessionTypes = []SessionType{SessionTypeText, SessionTypeImage}
)

// 客服聊天的会话记录
type CustomerSession struct {
	Id        string     `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 会话 ID
	Uid       string     `gorm:"not null;index;type:varchar(32)" json:"uid"`                   // 用户ID
	User      User       `gorm:"foreignkey:Uid" json:"user"`                                   // **外键**
	WaiterID  string     `gorm:"not null;index;type:varchar(32)" json:"waiter_id"`             // 客服 ID
	Waiter    User       `gorm:"foreignkey:WaiterID" json:"waiter"`                            //  **外键**
	ClosedAt  *time.Time `gorm:"null;index;" json:"closed_at"`                                 // 会话关闭时间
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (c *CustomerSession) TableName() string {
	return "customer_session"
}

func (c *CustomerSession) BeforeCreate(scope *gorm.Scope) error {
	// 不用自动生成，而是根据 md5(from + to) 生成
	//return scope.SetColumn("id", util.GenerateId())
	return nil
}

// 客服聊天的会话记录
type CustomerSessionItem struct {
	Id         string          `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 消息 ID
	SessionID  string          `gorm:"not null;index;type:varchar(32)" json:"session_id"`            // 会话 ID
	Session    CustomerSession `gorm:"foreignkey:SessionID" json:"session"`                          // **外键**
	Type       SessionType     `gorm:"not null;index;type:varchar(32)" json:"type"`                  // 会话类型
	SenderID   string          `gorm:"not null;index;type:varchar(32)" json:"sender_id"`             // 发送者 ID
	Sender     User            `gorm:"foreignkey:SenderID" json:"sender"`                            // **外键**
	ReceiverID string          `gorm:"not null;index;type:varchar(32)" json:"receiver_id"`           // 接受者的 ID
	Receiver   User            `gorm:"foreignkey:ReceiverID" json:"receiver"`                        // **外键**
	Payload    string          `gorm:"not null;index;type:text" json:"payload"`                      // 数据
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
}

func (c *CustomerSessionItem) TableName() string {
	return "customer_session_item"
}

func (c *CustomerSessionItem) IsValidType() bool {
	for _, t := range SessionTypes {
		if t == c.Type {
			return true
		}
	}

	return false
}

func (c *CustomerSessionItem) BeforeCreate(scope *gorm.Scope) error {
	if !c.IsValidType() {
		return exception.InvalidParams
	}
	return scope.SetColumn("id", util.GenerateId())
}
