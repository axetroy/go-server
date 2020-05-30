// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"regexp"
)

func LookUp(fullAreaCode string) (*Location, error) {
	var target Location

	reg := regexp.MustCompile("(\\d{2})(\\d{2})(\\d{2})")

	matcher := reg.FindStringSubmatch(fullAreaCode)

	if reg.MatchString(fullAreaCode) == false {
		return nil, exception.InvalidParams
	}

	provinceCode := matcher[1]
	cityCode := matcher[2]
	areaCode := matcher[3]

provinceLookup:
	for _, province := range Maps {
		if province.Code != provinceCode {
			continue provinceLookup
		}

		if province.FullCode == fullAreaCode {
			target = province
			break provinceLookup
		}

	cityLookup:
		for _, city := range province.Children {
			if city.Code != cityCode {
				continue cityLookup
			}

			if city.FullCode == fullAreaCode {
				target = city
				break
			}

		areaLookup:
			for _, area := range city.Children {
				if area.Code != areaCode {
					continue areaLookup
				}

				if area.FullCode == fullAreaCode {
					target = area
					break
				}
			}
		}
	}

	if target.Name == "" {
		return nil, exception.NoData
	}

	return &target, nil
}
