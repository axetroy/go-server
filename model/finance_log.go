package model

import "time"

type FinanceType string

var (
	FinanceTypeTransferIn  FinanceType = "transfer_in"  // 转入
	FinanceTypeTransferOut FinanceType = "transfer_out" // 转出

	FinanceLogMap = map[string]interface{}{
		"cny":  FinanceLogCny{},
		"usd":  FinanceLogUsd{},
		"coin": FinanceLogCoin{},
	}
)

type FinanceLog struct {
	Id       int64  `xorm:"pk unique notnull index" json:"id"` // 流水ID
	Currency string `xorm:"notnull index" json:"currency"`     // 对应的币种流水
	OrderId  int64  `xorm:"null index" json:"order_id"`        // 对应的订单id, 系统产生的流水可能不会orderId
	Uid      int64  `xorm:"notnull index" json:"uid"`          // 对应的用户

	BeforeBalance   float64 `xorm:"notnull" json:"before_balance"`   // 这条流水前的余额
	BalanceMutation float64 `xorm:"notnull" json:"balance_mutation"` // 可用余额的变动，正数则为加，负数为减
	AfterBalance    float64 `xorm:"notnull" json:"after_balance"`    // 这条流水后的余额

	BeforeFrozen   float64 `xorm:"notnull" json:"before_frozen"`   // 这条流水前的冻结余额
	FrozenMutation float64 `xorm:"notnull" json:"frozen_mutation"` // 冻结余额的变动,正数则为加，负数为减
	AfterFrozen    float64 `xorm:"notnull" json:"after_frozen"`    // 这条流水后的冻结余额

	Type FinanceType `xorm:"notnull" json:"status"` // 流水类型

	Note      *string   `xorm:"null varchar(128)" json:"note"` // 流水备注
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
}

type FinanceLogCny struct {
	FinanceLog `xorm:"extends"`
}

type FinanceLogUsd struct {
	FinanceLog `xorm:"extends"`
}

type FinanceLogCoin struct {
	FinanceLog `xorm:"extends"`
}
