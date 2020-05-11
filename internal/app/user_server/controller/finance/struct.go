// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package finance

import "github.com/axetroy/go-server/internal/model"

type Log struct {
	Id       string `json:"id"`       // 流水ID
	Currency string `json:"currency"` // 对应的币种流水
	OrderId  string `json:"order_id"` // 对应的订单id, 系统产生的流水可能不会orderId
	Uid      string `json:"uid"`      // 对应的用户

	BeforeBalance   float64 `json:"before_balance"`   // 这条流水前的余额
	BalanceMutation float64 `json:"balance_mutation"` // 可用余额的变动，正数则为加，负数为减
	AfterBalance    float64 `json:"after_balance"`    // 这条流水后的余额

	BeforeFrozen   float64 `json:"before_frozen"`   // 这条流水前的冻结余额
	FrozenMutation float64 `json:"frozen_mutation"` // 冻结余额的变动,正数则为加，负数为减
	AfterFrozen    float64 `json:"after_frozen"`    // 这条流水后的冻结余额

	Type model.FinanceType `json:"type"` // 流水类型

	Note      *string `json:"note"` // 流水备注
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
