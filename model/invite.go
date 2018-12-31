package model

import (
	"time"
)

type InviteStatus int

var (
	StatusInviteRegistered InviteStatus = 0  // 被邀请人刚注册
	StatusInviteAuth                    = 10 // 被邀请人进行了实名认证
	StatusInvitePay                     = 50 // 被邀请人已进行了一笔支付
)

type InviteHistory struct {
	Id            string       `xorm:"pk notnull unique index" json:"id"`
	Invitor       string       `xorm:"notnull index" json:"invitor"`        // 邀请人
	Invited       string       `xorm:"notnull unique index" json:"invited"` // 受邀请人, 只有唯一的一个
	Status        InviteStatus `xorm:"notnull" json:"status"`               // 受邀请人的激活状态
	RewardSettled bool         `xorm:"notnull" json:"reward_settled"`       // 是否已发放奖励, 包括邀请人和收邀请人的奖励
	CreatedAt     time.Time    `xorm:"created" json:"created_at"`
	UpdatedAt     time.Time    `xorm:"updated" json:"updated_at"`
	DeletedAt     *time.Time   `xorm:"deleted" json:"deleted_at"`
}
