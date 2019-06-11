// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news

import (
	"github.com/axetroy/go-server/exception"
)

var (
	ErrNewsInvalidType = exception.NewError("错误的文章类型")
	ErrNewsNotExist    = exception.NewError("文章不存在")
)
