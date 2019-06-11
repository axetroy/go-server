// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user_error

import (
	"github.com/axetroy/go-server/exception"
)

var (
	ErrUserExist                   = exception.NewError("用户已存在")
	ErrUserNotExist                = exception.NewError("用户不存在")
	ErrInvalidConfirmPassword      = exception.NewError("两次密码不一致")
	ErrInvalidResetCode            = exception.NewError("重置码错误或已失效")
	ErrErrRequireErrPayPasswordSet = exception.NewError("需要先设置交易密码")
)
