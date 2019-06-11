// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"github.com/axetroy/go-server/common_error"
)

var (
	// wallet
	ErrInvalidConfirmPassword      = common_error.NewError("两次密码不一致")
	ErrInvalidResetCode            = common_error.NewError("重置码错误或已失效")
	ErrErrRequireErrPayPasswordSet = common_error.NewError("需要先设置交易密码")
)
