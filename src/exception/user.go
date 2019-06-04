// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	// wallet
	InvalidConfirmPassword = New("两次密码不一致")
	InvalidResetCode       = New("重置码错误或已失效")
	RequirePayPasswordSet  = New("需要先设置交易密码")
)
