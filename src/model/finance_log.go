package model

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	"time"
)

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
	Id              string      `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"` // 流水ID
	Currency        string      `gorm:"not null;index;type:varchar(16)" json:"currency"`              // 对应的币种流水
	OrderId         string      `gorm:"null;index;type:varchar(32)" json:"order_id"`                  // 对应的订单id, 系统产生的流水可能不会存在orderId
	Uid             string      `gorm:"not null;index;type:varchar(32)" json:"uid"`                   // 对应的用户
	BeforeBalance   float64     `gorm:"not null;" json:"before_balance"`                              // 这条流水前的余额
	BalanceMutation float64     `gorm:"not null;" json:"balance_mutation"`                            // 可用余额的变动，正数则为加，负数为减
	AfterBalance    float64     `gorm:"not null" json:"after_balance"`                                // 这条流水后的余额
	BeforeFrozen    float64     `gorm:"not null" json:"before_frozen"`                                // 这条流水前的冻结余额
	FrozenMutation  float64     `gorm:"not null" json:"frozen_mutation"`                              // 冻结余额的变动,正数则为加，负数为减
	AfterFrozen     float64     `gorm:"not null" json:"after_frozen"`                                 // 这条流水后的冻结余额
	Type            FinanceType `gorm:"not null" json:"status"`                                       // 流水类型
	Note            *string     `gorm:"null;type:varchar(128)" json:"note"`                           // 流水备注
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       *time.Time `sql:"index" json:"-"`
}

type FinanceLogCny struct {
	FinanceLog
}

type FinanceLogUsd struct {
	FinanceLog
}

type FinanceLogCoin struct {
	FinanceLog
}

func (news *FinanceLog) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}

func (news *FinanceLogCny) TableName() string {
	return "finance_log_cny"
}

func (news *FinanceLogUsd) TableName() string {
	return "finance_log_usd"
}

func (news *FinanceLogCoin) TableName() string {
	return "finance_log_coin"
}

func (news *FinanceLogCny) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}

func (news *FinanceLogUsd) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}

func (news *FinanceLogCoin) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}
