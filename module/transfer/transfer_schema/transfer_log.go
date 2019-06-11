// Copyright 2019 Axetroy. All rights reserved. MIT license.
package transfer_schema

import (
	"github.com/axetroy/go-server/module/transfer/transfer_model"
)

type TransferLogPure struct {
	Id       string                        `json:"id"`       // 转账ID
	Currency string                        `json:"currency"` // 币种
	From     string                        `json:"from"`     // 谁转的
	To       string                        `json:"to"`       // 转给谁
	Amount   string                        `json:"amount"`   // 转账数量
	Status   transfer_model.TransferStatus `json:"status"`   // 转账状态
	Note     *string                       `json:"string"`   // 转账备注
}

type TransferLog struct {
	TransferLogPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
