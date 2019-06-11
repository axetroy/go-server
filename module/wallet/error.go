// Copyright 2019 Axetroy. All rights reserved. MIT license.
package wallet

import (
	"github.com/axetroy/go-server/common_error"
)

var (
	// wallet
	ErrNotEnoughBalance = common_error.NewError("钱包余额不足")
	ErrInvalidWallet    = common_error.NewError("无效的钱包")
)
