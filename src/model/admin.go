// Copyright 2019 Axetroy. All rights reserved. MIT license.
package model

import (
	"time"

	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
)

type AdminStatus int32

const (
	AdminStatusBanned      AdminStatus = -100 // 账号被禁用
	AdminStatusInactivated AdminStatus = -1   // 账号未激活
	AdminStatusInit        AdminStatus = 0    // 初始化状态
)

type Admin struct {
	Id        string      `gorm:"primary_key;not null;unique;index" json:"id"`            // 用户ID
	Username  string      `gorm:"not null;unique;index;type:varchar(36)" json:"username"` // 用户名, 用于登陆
	Name      string      `gorm:"not null;index;type:varchar(36)" json:"Name"`            // 管理员名
	Password  string      `gorm:"not null;type:varchar(36)" json:"password"`              // 登陆密码
	IsSuper   bool        `gorm:"not null;unique" json:"is_super"`                        // 是否是超级管理员, 超级管理员全站应该只有一个
	Status    AdminStatus `gorm:"not null;" json:"status"`                                // 状态
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (news *Admin) TableName() string {
	return "admin"
}

func (news *Admin) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}
