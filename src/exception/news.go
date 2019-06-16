// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	NewsInvalidType = New("错误的文章类型", 0)
	NewsNotExist    = New("文章不存在", 0)
)
