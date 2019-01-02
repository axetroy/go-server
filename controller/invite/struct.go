package invite

import (
	"github.com/axetroy/go-server/model"
)

type InvitePure struct {
	Id            string             `json:"id"`
	Invitor       string             `json:"invitor"`        // 邀请人
	Invited       string             `json:"invited"`        // 受邀请人, 只有唯一的一个
	Status        model.InviteStatus `json:"status"`         // 受邀请人的激活状态
	RewardSettled bool               `json:"reward_settled"` // 是否已发放奖励, 包括邀请人和收邀请人的奖励
}

type Invite struct {
	InvitePure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
