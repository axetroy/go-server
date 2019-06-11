// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner

import (
	"github.com/axetroy/go-server/common_error"
)

var (
	ErrBannerInvalidPlatform = common_error.NewError("无效的平台")
	ErrBannerNotExist        = common_error.NewError("不存在横幅")
)
