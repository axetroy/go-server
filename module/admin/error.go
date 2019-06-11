// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin

import (
	"github.com/axetroy/go-server/common_error"
)

var (
	ErrAdminExist    = common_error.NewError("管理员已存在")
	ErrAdminNotExist = common_error.NewError("管理员不存在")
	ErrAdminNotSuper = common_error.NewError("只有超级管理员才能操作")
)
