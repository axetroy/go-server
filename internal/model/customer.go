// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.

package model

import (
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/jinzhu/gorm"
	"time"
)

type Customer struct {
	Id        string `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 客服 ID
	Username  string `gorm:"not null;unique;index;type:varchar(36)" json:"username"`       // 用户名, 用于登陆
	NickName  string `gorm:"not null;index;type:varchar(36)" json:"nick_name"`             // 昵称
	Password  string `gorm:"not null;type:varchar(36)" json:"password"`                    // 登陆密码
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (c *Customer) TableName() string {
	return "customer"
}

func (c *Customer) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}
