// Copyright 2019 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/core/util"
	"github.com/jinzhu/gorm"
	"time"
)

type LoginLogType int
type LoginLogCommand int

const (
	LoginLogTypeUserName         LoginLogType    = 0 // 用户名登陆
	LoginLogTypeTel                                  // 手机登陆
	LoginLogTypeEmail                                // 邮箱登陆
	LoginLogTypeThird                                // 第三方登陆
	LoginLogTypeWechat                               // 微信登陆
	LoginLogCommandLoginSuccess  LoginLogCommand = 0 // 登陆成功
	LoginLogCommandLogoutSuccess                     // 登出成功
	LoginLogCommandLoginFail                         // 登陆失败
	LoginLogCommandLogoutFail                        // 登出失败
)

type LoginLog struct {
	Id        string          `gorm:"primary_key;not null;index;type:varchar(32)" json:"id"` // 数据ID
	Uid       string          `gorm:"not null;index;type:varchar(32)" json:"uid"`            // 用户ID
	User      User            `gorm:"foreignkey:Uid" json:"user"`                            // **外键**
	Type      LoginLogType    `gorm:"not null;type:int" json:"type"`                         // 登陆类型(用什么方式登陆)
	Command   LoginLogCommand `gorm:"not null;type:int" json:"command"`                      // 登陆的状态(成功, 失败)
	LastIp    string          `gorm:"not null;type:varchar(15)" json:"last_ip"`              // 本次登陆IP
	Client    string          `gorm:"not null;type:varchar(255)" json:"client"`              // 登陆的客户端
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (news *LoginLog) TableName() string {
	return "login_log"
}

func (news *LoginLog) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}
