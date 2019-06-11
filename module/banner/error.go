// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner

import (
	"github.com/axetroy/go-server/exception"
)

var (
	ErrBannerInvalidPlatform = exception.NewError("无效的平台")
	ErrBannerNotExist        = exception.NewError("不存在横幅")
)
