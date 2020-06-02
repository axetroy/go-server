// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.

package model

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/jinzhu/gorm"
	"time"
)

type SessionType string

const (
	SessionTypeDistribution SessionType = "distribution" // 分配客服，客户端收到这个时间，说明已经连接到了客服
	SessionTypeText         SessionType = "text"         // 发送纯文本消息
	SessionTypeImage        SessionType = "image"        // 发送图片
	SessionTypeTip          SessionType = "tip"          // 系统发送提示
)

var (
	SessionTypes = []SessionType{SessionTypeText, SessionTypeImage, SessionTypeTip}
)

// 客服聊天的会话记录
type CustomerSession struct {
	Id         string   `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 会话 ID
	Uid        string   `gorm:"not null;index;type:varchar(32)" json:"uid"`                   // 用户ID
	User       User     `gorm:"foreignkey:Uid" json:"user"`                                   // **外键**
	CustomerID string   `gorm:"not null;index;type:varchar(32)" json:"customer_id"`           // 客服 ID
	Customer   Customer `gorm:"foreignkey:CustomerID" json:"customer_id"`                     //  **外键**
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeletedAt  *time.Time `sql:"index"`
}

func (c *CustomerSession) TableName() string {
	return "customer_session"
}

func (c *CustomerSession) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}

// 客服聊天的会话记录
type CustomerSessionItem struct {
	Id        string      `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 会话 ID
	SessionId string      `gorm:"not null;index;type:varchar(32)" json:"user_id"`               // 用户 ID
	Type      SessionType `gorm:"not null;index;type:varchar(32)" json:"type"`                  // 会话类型
	Payload   string      `gorm:"not null;index;type:text" json:"payload"`                      // 数据
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
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
