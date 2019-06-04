// Copyright 2019 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/src/util"
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
	LoginLogCommandLoginSuccess  LoginLogCommand = 0 // 登陆成功
	LoginLogCommandLogoutSuccess                     // 登出成功
	LoginLogCommandLoginFail                         // 登陆失败
	LoginLogCommandLogoutFail                        // 登出失败
)

type LoginLog struct {
	Id        string          `gorm:"primary_key;not null;index;type:varchar(32)" json:"id"`
	Uid       string          `gorm:"not null;index;type:varchar(32)" json:"uid"`
	Type      LoginLogType    `gorm:"not null;type:int" json:"type"`
	Command   LoginLogCommand `gorm:"not null;type:int" json:"command"`
	LastIp    string          `gorm:"not null;type:varchar(15)" json:"last_ip"`
	Client    string          `gorm:"not null;type:varchar(255)" json:"client"`
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
