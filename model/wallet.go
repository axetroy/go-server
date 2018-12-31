package model

import "time"

var (
	WalletCNY  = "CNY"
	WalletUSD  = "USD"
	WalletCOIN = "COIN"
)

type Wallet struct {
	Id        int64     `xorm:"pk unique notnull index" json:"-"`
	Balance   float64   `json:"balance"`
	Frozen    float64   `json:"frozen"`
	CreatedAt time.Time `xorm:"created" json:"created_at"`
	UpdatedAt time.Time `xorm:"updated" json:"updated_at"`
	DeletedAt time.Time `xorm:"deleted" json:"-"`
}

// 人民币
type WalletCny struct {
	Wallet `xorm:"extends"`
}

// 美元
type WalletUsd struct {
	Wallet `xorm:"extends"`
}

// 我们平台自己的币
type WalletCoin struct {
	Wallet `xorm:"extends"`
}
