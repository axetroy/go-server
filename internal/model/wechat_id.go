// Copyright 2019 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/jinzhu/gorm"
	"time"
)

type WechatOpenID struct {
	Id   string `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 微信的 open ID
	Uid  string `gorm:"null;unique;index" json:"uid"`                                 // 对应的用户 ID, 如果为空，则说明没有关联帐号
	User User   `gorm:"foreignkey:Uid" json:"user"`                                   // **外键**

	// 微信相关字段 https://developers.weixin.qq.com/miniprogram/dev/api/open-api/user-info/UserInfo.html
	Nickname  *string `gorm:"null;" json:"nickname"`   // 微信昵称
	AvatarUrl *string `gorm:"null;" json:"avatar_url"` // 微信头像 URL
	Gender    *int    `gorm:"null;" json:"gender"`     // 微信性别 0: 未知 1: 男性 2: 女性
	Country   *string `gorm:"null;" json:"country"`    // 国家/地区
	Province  *string `gorm:"null;" json:"province"`   // 省份
	City      *string `gorm:"null;" json:"city"`       // 城市
	Language  *string `gorm:"null;" json:"language"`   // 语言

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
