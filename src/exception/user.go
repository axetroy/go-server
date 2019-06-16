// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	UserExist             = New("用户已存在", 0)
	UserNotExist          = New("用户不存在", 0)
	InvalidResetCode      = New("重置码错误或已失效", 0)
	RequirePayPasswordSet = New("需要先设置交易密码", 0)
)
