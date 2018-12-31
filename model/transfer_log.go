package model

import "time"

type TransferStatus int

var (
	TransferStatusReject         TransferStatus = -1 // 收款方拒接接受
	TransferStatusWaitForConfirm TransferStatus = 0  // 等待收款方确认
	TransferStatusConfirmed      TransferStatus = 0  // 收款方已确认

	TransferLogMap = map[string]interface{}{
		"cny":  TransferLogCny{},
		"usd":  TransferLogUsd{},
		"coin": TransferLogCoin{},
	}
)

type TransferLog struct {
	Id           int64          `xorm:"pk unique notnull index" json:"id"` // 转账ID
	Currency     string         `xorm:"notnull" json:"currency"`           // 转账币种
	From         int64          `xorm:"notnull index" json:"from"`         // 汇款人
	To           int64          `xorm:"notnull index" json:"to"`           // 收款人
	Amount       float64        `xorm:"notnull" json:"amount"`             // 转账数量
	Status       TransferStatus `xorm:"notnull" json:"status"`             // 转账状态
	Note         *string        `xorm:"null varchar(128)" json:"note"`     // 转账备注
	SnapshotFrom *string        `xorm:"null ->" json:"-"`                  // 转账者的钱包快照
	SnapshotTo   *string        `xorm:"null ->" json:"-"`                  // 收款人的钱包快照
	CreatedAt    time.Time      `xorm:"created" json:"created_at"`
	UpdatedAt    time.Time      `xorm:"updated" json:"updated_at"`
	DeletedAt    time.Time      `xorm:"deleted" json:"-"`
}

type TransferLogCny struct {
	TransferLog `xorm:"extends"`
}

type TransferLogUsd struct {
	TransferLog `xorm:"extends"`
}

type TransferLogCoin struct {
	TransferLog `xorm:"extends"`
}
