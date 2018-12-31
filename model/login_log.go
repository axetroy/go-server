package model

import "time"

type LoginLog struct {
	Id       int64  `json:"id"`
	Uid      string `json:"uid"`
	Username string `json:"username"`
	// 登录方式
	// 0用户名 	1手机	2邮箱	3第三方
	Type int32 `xorm:"notnull" json:"type"`
	// 操作类型
	// 1登陆成功  2登出成功 3登录失败 4登出失败
	Command   int32     `xorm:"notnull" json:"command"`
	LastIp    string    `xorm:"notnull" json:"last_ip"`
	Client    string    `xorm:"notnull" json:"client"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"deleted_at"`
	Version   int32     `xorm:"version" json:"version"`
}
