// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin

import (
	"github.com/axetroy/go-server/exception"
)

var (
	ErrAdminExist    = exception.NewError("管理员已存在")
	ErrAdminNotExist = exception.NewError("管理员不存在")
	ErrAdminNotSuper = exception.NewError("只有超级管理员才能操作")
)
