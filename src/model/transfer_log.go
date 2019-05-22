package model

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

type TransferStatus int

var (
	transferLogTablePrefix                      = "transfer_log_"
	TransferStatusReject         TransferStatus = -1 // 收款方拒接接受
	TransferStatusWaitForConfirm TransferStatus = 0  // 等待收款方确认
	TransferStatusConfirmed      TransferStatus = 1  // 收款方已确认

	TransferLogCnyTableName  = transferLogTablePrefix + strings.ToLower(WalletCNY)  // 人民币表名
	TransferLogUsdTableName  = transferLogTablePrefix + strings.ToLower(WalletUSD)  // 美元表名
	TransferLogCoinTableName = transferLogTablePrefix + strings.ToLower(WalletCOIN) // 积分表名

	TransferTableNames = []string{ // 所有的表名
		TransferLogCnyTableName,
		TransferLogUsdTableName,
		TransferLogCoinTableName,
	}

	TransferLogMap = map[string]interface{}{
		WalletCNY:  TransferLogCny{},
		WalletUSD:  TransferLogUsd{},
		WalletCOIN: TransferLogCoin{},
	}
)

type TransferLog struct {
	Id           string         `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"` // 转账ID
	Currency     string         `gorm:"not null;index;" json:"currency"`                              // 转账币种
	From         string         `gorm:"not null;index;type:varchar(32)" json:"from"`                  // 汇款人
	To           string         `gorm:"not null;index;type:varchar(32)" json:"to"`                    // 收款人
	Amount       string         `gorm:"not null;type:numeric" json:"amount"`                          // 转账数量
	Status       TransferStatus `gorm:"not null" json:"status"`                                       // 转账状态
	Note         *string        `gorm:"null;type:varchar(128)" json:"note"`                           // 转账备注
	SnapshotFrom *string        `gorm:"null" json:"-"`                                                // 转账者的钱包快照
	SnapshotTo   *string        `gorm:"null" json:"-"`                                                // 收款人的钱包快照
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time `sql:"index" json:"-"`
}

type TransferLogCny struct {
	TransferLog
}

type TransferLogUsd struct {
	TransferLog
}

type TransferLogCoin struct {
	TransferLog
}

func (news *TransferLog) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}

func (news *TransferLogCny) TableName() string {
	return TransferLogCnyTableName
}

func (news *TransferLogUsd) TableName() string {
	return TransferLogUsdTableName
}

func (news *TransferLogCoin) TableName() string {
	return TransferLogCoinTableName
}
