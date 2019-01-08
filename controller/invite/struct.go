package invite

import (
	"github.com/axetroy/go-server/model"
)

type Pure struct {
	Id            string             `json:"id"`
	Inviter       string             `json:"inviter"`        // 邀请人
	Invitee       string             `json:"invitee"`        // 受邀请人, 只有唯一的一个
	Status        model.InviteStatus `json:"status"`         // 受邀请人的激活状态
	RewardSettled bool               `json:"reward_settled"` // 是否已发放奖励, 包括邀请人和收邀请人的奖励
}

type Invite struct {
	Pure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
