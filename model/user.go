package model

import (
	"time"
)

type UserStatus int32

type Gender int

const (
	// 用户状态
	UserStatusBanned      UserStatus = -100 // 账号被禁用
	UserStatusInactivated            = -1   // 账号未激活
	UserStatusInit                   = 0    // 初始化状态

	// 用户性别
	GenderUnknown Gender = 0 // 未知性别
	GenderMale               // 男
	GenderFemmale            // 女
)

type User struct {
	Id          int64      `xorm:"pk notnull unique index" json:"id"`            // 用户ID
	Username    string     `xorm:"notnull unique index" json:"username"`         // 用户名
	Password    string     `xorm:"notnull varchar(36)" json:"password"`          // 登陆密码
	PayPassword *string    `xorm:"varchar(36)" json:"pay_password"`              // 支付密码
	Nickname    *string    `xorm:"null varchar(36)" json:"nickname"`             // 昵称
	Phone       *string    `xorm:"null varchar(16) index" json:"phone"`          // 手机号
	Email       *string    `xorm:"null varchar(36) index" json:"email"`          // 邮箱
	Status      UserStatus `xorm:"notnull" json:"status"`                        // 状态
	Role        string     `xorm:"notnull varchar(36)" json:"role"`              // 角色
	Avatar      string     `xorm:"notnull varchar(36)" json:"avatar"`            // 头像
	Level       int32      `xorm:"default(1)" json:"level"`                      // 用户等级
	Gender      Gender     `xorm:"default(0)" json:"gender"`                     // 性别
	InviteCode  string     `xorm:"notnull unique varchar(8)" json:"invite_code"` // 邀请码
	CreatedAt   time.Time  `xorm:"created" json:"created_at"`
	UpdatedAt   time.Time  `xorm:"updated" json:"updated_at"`
	DeletedAt   *time.Time `xorm:"deleted" json:"deleted_at"`
}
