// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

func IsValid(provinceCode string, cityCode string, areaCode string) bool {
	for _, province := range Maps {
		if province.FullCode == provinceCode {
			for _, city := range province.Children {
				if city.FullCode == cityCode {
					for _, area := range city.Children {
						if area.FullCode == areaCode {
							return true
						}
					}
				}
			}
		}
	}

	return false
}
