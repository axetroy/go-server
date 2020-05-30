// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"github.com/axetroy/go-server/internal/library/exception"
)

func LookUp(fullAreaCode string) (*Location, error) {
	var target Location

	// 查找省份
	if len(fullAreaCode) == 2 {
		if name, ok := ProvinceMap[fullAreaCode]; ok {
			target.Code = fullAreaCode
			target.Name = name
			return &target, nil
		} else {
			return nil, exception.NoData
		}
	}

	// 查找城市
	if len(fullAreaCode) == 4 {
		if name, ok := CityMap[fullAreaCode]; ok {
			target.Code = fullAreaCode
			target.Name = name
			return &target, nil
		} else {
			return nil, exception.NoData
		}
	}

	// 查找地区
	if len(fullAreaCode) == 6 {
		if name, ok := AreaMap[fullAreaCode]; ok {
			target.Code = fullAreaCode
			target.Name = name
			return &target, nil
		} else {
			return nil, exception.NoData
		}
	}

	// 查找地区
	if len(fullAreaCode) == 9 {
		if name, ok := StreetMap[fullAreaCode]; ok {
			target.Code = fullAreaCode
			target.Name = name
			return &target, nil
		} else {
			return nil, exception.NoData
		}
	}

	return nil, exception.NoData
}
