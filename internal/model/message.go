// Copyright 2019 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/internal/util"
	"github.com/jinzhu/gorm"
	"time"
)

type MessageStatus int

const (
	MessageStatusActive MessageStatus = 0 // 默认状态
)

type Message struct {
	Id        string        `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 消息ID
	Uid       string        `gorm:"not null;index;type:varchar(32)" json:"uid"`                   // 这条消息的所有者
	Title     string        `gorm:"not null;index;type:varchar(32)" json:"title"`                 // 消息标题
	Content   string        `gorm:"not null;type:text" json:"content"`                            // 消息内容
	Read      bool          `gorm:"not null" json:"read"`                                         // 是否已读
	ReadAt    *time.Time    `json:"read_at"`                                                      // 已读时间
	Status    MessageStatus `gorm:"not null" json:"status"`                                       // 消息状态
	Note      *string       `gorm:"null;type:varchar(255)" json:"note"`                           // 这条通知的备注
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (news *Message) TableName() string {
	return "message"
}

func (news *Message) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", util.GenerateId()); err != nil {
		return err
	}
	return nil
}
