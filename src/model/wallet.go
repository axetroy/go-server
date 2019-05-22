package model

import (
	"strings"
	"time"
)

var (
	walletTablePrefix = "wallet_"
	WalletCNY         = "CNY"
	WalletUSD         = "USD"
	WalletCOIN        = "COIN"
	Wallets           = []string{WalletCNY, WalletUSD, WalletCOIN}

	WalletCnyTableName  = walletTablePrefix + strings.ToLower(WalletCNY)  // 人民币表名
	WalletUsdTableName  = walletTablePrefix + strings.ToLower(WalletUSD)  // 美元表名
	WalletCoinTableName = walletTablePrefix + strings.ToLower(WalletCOIN) // 积分表名

	//WalletTableNames = []string{ // 所有的表名
	//	WalletCnyTableName,
	//	WalletUsdTableName,
	//	WalletCoinTableName,
	//}
)

type Wallet struct {
	Id        string  `gorm:"primary_key;unique;notnull;index;type:varchar(32)" json:"id"` // 用户ID
	Currency  string  `gorm:"not null;type:varchar(12)" json:"currency"`                   // 钱包币种
	Balance   float64 `gorm:"not null;type:numeric" json:"balance"`                        // 可用余额
	Frozen    float64 `gorm:"not null;type:numeric" json:"frozen"`                         // 冻结余额
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

// 人民币
type WalletCny struct {
	Wallet
}

// 美元
type WalletUsd struct {
	Wallet
}

// 我们平台自己的币
type WalletCoin struct {
	Wallet
}

func (news *WalletCny) TableName() string {
	return WalletCnyTableName
}

func (news *WalletUsd) TableName() string {
	return WalletUsdTableName
}

func (news *WalletCoin) TableName() string {
	return WalletCoinTableName
}
