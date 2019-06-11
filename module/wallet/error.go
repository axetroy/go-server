// Copyright 2019 Axetroy. All rights reserved. MIT license.
package wallet

import (
	"github.com/axetroy/go-server/exception"
)

var (
	// wallet
	ErrNotEnoughBalance = exception.NewError("钱包余额不足")
	ErrInvalidWallet    = exception.NewError("无效的钱包")
)
