// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	// wallet
	NotEnoughBalance = New("钱包余额不足")
	InvalidWallet    = New("无效的钱包")
)
