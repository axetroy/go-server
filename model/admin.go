package model

import (
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
