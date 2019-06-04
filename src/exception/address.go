// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	AddressDefaultNotExist     = New("默认地址不存在")
	AddressNotExist            = New("地址记录不存在")
	AddressInvalidProvinceCode = New("无效的省份代码")
	AddressInvalidCityCode     = New("无效的城市代码")
	AddressInvalidAreaCode     = New("无效的地区代码")
)
