// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news

import (
	"github.com/axetroy/go-server/common_error"
)

var (
	ErrNewsInvalidType = common_error.NewError("错误的文章类型")
	ErrNewsNotExist    = common_error.NewError("文章不存在")
)
