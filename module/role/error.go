// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"github.com/axetroy/go-server/exception"
)

var (
	ErrRoleNotExist     = exception.NewError("角色不存在")
	ErrRoleCannotUpdate = exception.NewError("无法更新角色")
	ErrRoleHadBeenUsed  = exception.NewError("角色正在被使用，无法删除")
)
