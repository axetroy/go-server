// Copyright 2019 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	"time"
)

type BannerPlatform string

const (
	BannerPlatformPc  BannerPlatform = "pc"  // PC 端的 Banner
	BannerPlatformApp BannerPlatform = "app" // APP 端的 Banner
)

type Banner struct {
	w
	Id          string         `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"` // ID
	Image       string         `gorm:"not null;index;type:varchar(255)" json:"image"`                // 图片
	Href        string         `gorm:"not null;index;type:varchar(255)" json:"href"`                 // 图片连接
	Platform    BannerPlatform `gorm:"not null;index;type:varchar(32)" json:"platform"`              // 用于哪个平台
	Description *string        `gorm:"null;index;type:varchar(255)" json:"description"`              // Banner 描述
	Priority    *int           `gorm:"null;index;" json:"priority"`                                  // 优先级，主要用于排序
	Identifier  *string        `gorm:"null;index;type:varchar(32)" json:"identifier"`                // 标识符, 用于 APP 跳转页面的标识符
	FallbackUrl *string        `gorm:"null;index;type:varchar(255)" json:"fallback_url"`             // fallback 的 url， 当 APP 没有 `Identifier` 对应的页面时，这个就是 fallback 的页面
	Active      bool           `gorm:"not null;default:true;index;" json:"active"`                   // 是否激活
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

func (news *Banner) TableName() string {
	return "banner"
}

func (news *Banner) BeforeCreate(scope *gorm.Scope) (err error) {
	if err = scope.SetColumn("id", util.GenerateId()); err != nil {
		return
	}

	if err = scope.SetColumn("active", true); err != nil {
		return
	}

	return
}
