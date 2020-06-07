// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import "strings"

func IsValid(provinceCode string, cityCode string, areaCode string, streetCode string) bool {
	if len(provinceCode) != 2 || len(cityCode) != 4 || len(areaCode) != 6 || len(streetCode) != 9 {
		return false
	}

	if !strings.HasPrefix(cityCode, provinceCode) {
		return false
	}

	if !strings.HasPrefix(areaCode, cityCode) {
		return false
	}

	if !strings.HasPrefix(areaCode, cityCode) {
		return false
	}

	if !strings.HasPrefix(streetCode, areaCode) {
		return false
	}

	if _, ok := StreetMap[streetCode]; ok {
		return true
	}

	return false
}
