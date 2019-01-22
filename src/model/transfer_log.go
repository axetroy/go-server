package model

import "time"

type TransferStatus int

var (
	TransferStatusReject         TransferStatus = -1 // 收款方拒接接受
	TransferStatusWaitForConfirm TransferStatus = 0  // 等待收款方确认
	TransferStatusConfirmed      TransferStatus = 1  // 收款方已确认

	TransferLogMap = map[string]interface{}{
		"cny":  TransferLogCny{},
		"usd":  TransferLogUsd{},
		"coin": TransferLogCoin{},
	}
)

type TransferLog struct {
	Id           string         `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"` // 转账ID
	Currency     string         `gorm:"not null;index;" json:"currency"`                              // 转账币种
	From         string         `gorm:"not null;index;type:varchar(32)" json:"from"`                  // 汇款人
	To           string         `gorm:"not null;index;type:varchar(32)" json:"to"`                    // 收款人
	Amount       float64        `gorm:"not null;" json:"amount"`                                      // 转账数量
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

func (news *TransferLogCny) TableName() string {
	return "transfer_log_cny"
}

func (news *TransferLogUsd) TableName() string {
	return "transfer_log_usd"
}

func (news *TransferLogCoin) TableName() string {
	return "transfer_log_coin"
}
