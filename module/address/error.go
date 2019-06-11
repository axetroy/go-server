// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address

import (
	"github.com/axetroy/go-server/exception"
)

var (
	ErrDefaultAddressNotExist     = exception.NewError("默认地址不存在")
	ErrAddressNotExist            = exception.NewError("地址记录不存在")
	ErrAddressInvalidProvinceCode = exception.NewError("无效的省份代码")
	ErrAddressInvalidCityCode     = exception.NewError("无效的城市代码")
	ErrAddressInvalidAreaCode     = exception.NewError("无效的地区代码")
)
