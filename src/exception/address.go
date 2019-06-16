// Copyright 2019 Axetroy. All rights reserved. MIT license.
package exception

var (
	AddressDefaultNotExist     = New("默认地址不存在", 0)
	AddressNotExist            = New("地址记录不存在", 0)
	AddressInvalidProvinceCode = New("无效的省份代码", 0)
	AddressInvalidCityCode     = New("无效的城市代码", 0)
	AddressInvalidAreaCode     = New("无效的地区代码", 0)
)
