// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package schema

type AddressPure struct {
	Id           string `json:"id"`            // 地址ID
	Name         string `json:"name"`          // 收货人
	Phone        string `json:"phone"`         // 收货人手机号
	ProvinceCode string `json:"province_code"` // 省份代码
	CityCode     string `json:"city_code"`     // 城市代码
	AreaCode     string `json:"area_code"`     // 区域代码
	StreetCode   string `json:"street_code"`   // 街道/乡镇代码
	Address      string `json:"address"`       // 详细的地址
	IsDefault    bool   `json:"is_default"`    // 是否是默认地址
}

type Address struct {
	AddressPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
