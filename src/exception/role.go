// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	RoleNotExist     = New("角色不存在", 0)
	RoleCannotUpdate = New("无法更新角色", 0)
	RoleHadBeenUsed  = New("角色正在被使用，无法删除", 0)
)
