// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"github.com/axetroy/go-server/common_error"
)

var (
	ErrRoleNotExist     = common_error.NewError("角色不存在")
	ErrRoleCannotUpdate = common_error.NewError("无法更新角色")
	ErrRoleHadBeenUsed  = common_error.NewError("角色正在被使用，无法删除")
)
