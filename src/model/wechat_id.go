// Copyright 2019 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type WechatOpenID struct {
	Id        string `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 微信的 open ID
	Uid       string `gorm:"not null;unique;index" json:"uid"`                             // 对应的用户 ID
	User      User   `gorm:"foreignkey:Uid" json:"user"`                                   // **外键**
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (news *WechatOpenID) TableName() string {
	return "wechat_open_id"
}

func (news *WechatOpenID) BeforeCreate(scope *gorm.Scope) error {
	return nil
}
