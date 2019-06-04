// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	RequireFile    = New("请上传文件")
	NotSupportType = New("不支持该文件类型")
	OutOfSize      = New("超出文件大小限制")
)
