package model

import (
	"time"
)

type InviteType string

var (
	InviteTypeInitiative InviteType = "initiative" // 类型为邀请人
	InviteTypePassive    InviteType = "passive"    // 类型为被邀请人
)

type InviteRewardHistory struct {
	Id        string     `xorm:"pk notnull unique" json:"id"`
	InviteId  string     `xorm:"notnull" json:"invitor"` // 邀请记录ID
	Uid       string     `xorm:"notnull" json:"invited"` // 收益人ID
	Type      InviteType `xorm:"notnull" json:"type"`    // 邀请类型
	CreatedAt time.Time  `xorm:"created" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated" json:"updated_at"`
	DeletedAt time.Time  `xorm:"deleted" json:"deleted_at"`
}
