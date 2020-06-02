// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package status

import (
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
)

type Status struct {
	OnlineWaiterNum int                 `json:"online_waiter_num"` // 当前在线的客服
	OnlineUserNum   int                 `json:"online_user_num"`   // 当前在线的用户
	PendingNum      int                 `json:"pending_num"`       // 正在等待的用户数量
	Matchers        map[string][]string `json:"matchers"`          // 当前正在配对的连接
}

var GetStatusRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return schema.Response{
			Message: "",
			Status:  schema.StatusSuccess,
			Data: Status{
				OnlineUserNum:   ws.UserPoll.Length(),
				OnlineWaiterNum: ws.WaiterPoll.Length(),
				PendingNum:      ws.MatcherPool.GetPendingLength(),
				Matchers:        ws.MatcherPool.GetMatcher(),
			},
		}
	})
})
