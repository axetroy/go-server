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
	Id        int64      `xorm:"pk notnull unique" json:"id"`
	InviteId  int64      `xorm:"notnull" json:"invitor"` // 邀请记录ID
	Uid       int64      `xorm:"notnull" json:"invited"` // 收益人ID
	Type      InviteType `xorm:"notnull" json:"invited"` // 邀请类型
	CreatedAt time.Time  `xorm:"created" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated" json:"updated_at"`
	DeletedAt time.Time  `xorm:"deleted" json:"deleted_at"`
}
