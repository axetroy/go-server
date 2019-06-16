// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	AdminExist    = New("管理员已存在", 0)
	AdminNotExist = New("管理员不存在", 0)
	AdminNotSuper = New("只有超级管理员才能操作", 0)
)
