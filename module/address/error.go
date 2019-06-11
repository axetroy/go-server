// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address

import (
	"github.com/axetroy/go-server/common_error"
)

var (
	ErrDefaultAddressNotExist     = common_error.NewError("默认地址不存在")
	ErrAddressNotExist            = common_error.NewError("地址记录不存在")
	ErrAddressInvalidProvinceCode = common_error.NewError("无效的省份代码")
	ErrAddressInvalidCityCode     = common_error.NewError("无效的城市代码")
	ErrAddressInvalidAreaCode     = common_error.NewError("无效的地区代码")
)
