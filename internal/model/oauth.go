// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package model

import (
	"time"

	"github.com/axetroy/go-server/internal/library/util"
	"github.com/jinzhu/gorm"
)

type OAuthProvider string

var (
	ProviderGithub   OAuthProvider = "github"
	ProviderGitlab   OAuthProvider = "gitlab"
	ProviderTwitter  OAuthProvider = "twitter"
	ProviderFacebook OAuthProvider = "facebook"
	ProviderGoogle   OAuthProvider = "google"
	providerMap                    = map[OAuthProvider]bool{
		ProviderGithub:   true,
		ProviderGitlab:   true,
		ProviderTwitter:  true,
		ProviderFacebook: true,
		ProviderGoogle:   true,
	}
)

type OAuth struct {
	Id       string        `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"`    // ID
	Provider OAuthProvider `gorm:"not null;index;unique(user_id);type:varchar(32)" json:"provider"` // 绑定的服务提供商
	Uid      string        `gorm:"not null;index;unique(provider);type:varchar(32)" json:"uid"`     // 对应的平台 UID

	// 以下是 oAuth 返回的字段
	UserID            string    `gorm:"not null;type:varchar(255);index" json:"user_id"` // 用户 ID
	Name              string    `gorm:"not null;type:varchar(32)" json:"name"`           // 用户名
	FirstName         string    `gorm:"not null;type:varchar(32)" json:"first_name"`     // 姓
	LastName          string    `gorm:"not null;type:varchar(32)" json:"last_name"`      // 名
	Nickname          string    `gorm:"not null;type:varchar(32)" json:"nickname"`       // 昵称
	Description       string    `gorm:"not null;type:varchar(255)" json:"description"`   // 描述
	Email             string    `gorm:"not null;type:varchar(255)" json:"email"`         // 邮箱
	AvatarURL         string    `gorm:"not null;type:varchar(255)" json:"avatar_url"`    // 头像 URL
	Location          string    `gorm:"not null;type:varchar(255)" json:"location"`      // 地址
	AccessToken       string    `gorm:"not null;" json:"access_token"`                   // Access Token
	AccessTokenSecret string    `gorm:"not null;" json:"access_token_secret"`            // Access Token Secret
	RefreshToken      string    `gorm:"not null;" json:"refresh_token"`                  // Refresh Token
	ExpiresAt         time.Time `gorm:"not null;" json:"expires_at"`                     // 过期时间

	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (o *OAuth) TableName() string {
	return "oauth"
}

func (o *OAuth) IsValidProvider() bool {
	if _, ok := providerMap[o.Provider]; ok == false {
		return false
	}
	return true
}

func (o *OAuth) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("id", util.GenerateId())
	return
}
