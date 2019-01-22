package model

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	"time"
)

type InviteStatus int

var (
	StatusInviteRegistered InviteStatus = 0  // 被邀请人刚注册
	StatusInviteAuth                    = 10 // 被邀请人进行了实名认证
	StatusInvitePay                     = 50 // 被邀请人已进行了一笔支付
)

type InviteHistory struct {
	Id            string       `gorm:"primary_key;notnull;unique;index" json:"id"`
	Inviter       string       `gorm:"not null;index;type:varchar(32)" json:"inviter"`        // 邀请人
	Invitee       string       `gorm:"not null;unique;index;type:varchar(32)" json:"invitee"` // 受邀请人, 只有唯一的一个
	Status        InviteStatus `gorm:"not null;" json:"status"`                               // 受邀请人的激活状态
	RewardSettled bool         `gorm:"not null;" json:"reward_settled"`                       // 是否已发放奖励, 包括邀请人和收邀请人的奖励
	CreatedAt     time.Time
	UpdatedAt     time.Time
	DeletedAt     *time.Time `sql:"index"`
}

func (news *InviteHistory) TableName() string {
	return "invite_history"
}

func (news *InviteHistory) BeforeCreate(scope *gorm.Scope) error {
	// 生成ID
	if err := scope.SetColumn("id", util.GenerateId()); err != nil {
		return err
	}
	return nil
}
