package model

import (
	"github.com/axetroy/go-server/id"
	"github.com/jinzhu/gorm"
	"time"
)

type AdminStatus int32

const (
	// 管理员状态
	AdminStatusBanned      AdminStatus = -100 // 账号被禁用
	AdminStatusInactivated             = -1   // 账号未激活
	AdminStatusInit                    = 0    // 初始化状态
)

type Admin struct {
	Id        string      `xorm:"pk notnull unique index" json:"id"`                // 用户ID
	Username  string      `xorm:"notnull unique index varchar(36)" json:"username"` // 用户名, 用于登陆
	Name      string      `xorm:"notnull index varchar(36)" json:"Name"`            // 管理员名
	Password  string      `xorm:"notnull varchar(36)" json:"password"`              // 登陆密码
	IsSuper   bool        `xorm:"notnull unique" json:"is_super"`                   // 是否是超级管理员, 超级管理员全站应该只有一个
	Status    AdminStatus `xorm:"notnull" json:"status"`                            // 状态
	CreatedAt time.Time   `xorm:"created" json:"created_at"`
	UpdatedAt time.Time   `xorm:"updated" json:"updated_at"`
	DeletedAt *time.Time  `xorm:"deleted" json:"deleted_at"`
}

type AdminGo struct {
	Id        string      `gorm:"primary_key;not null;unique;index" json:"id"`            // 用户ID
	Username  string      `gorm:"not null;unique;index;type:varchar(36)" json:"username"` // 用户名, 用于登陆
	Name      string      `gorm:"not null;indextype:varchar(36)" json:"Name"`             // 管理员名
	Password  string      `gorm:"not null;type:varchar(36)" json:"password"`              // 登陆密码
	IsSuper   bool        `gorm:"not null;unique" json:"is_super"`                        // 是否是超级管理员, 超级管理员全站应该只有一个
	Status    AdminStatus `gorm:"not null;" json:"status"`                                // 状态
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (news *AdminGo) TableName() string {
	return "admin"
}

func (news *AdminGo) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", id.Generate())
}
