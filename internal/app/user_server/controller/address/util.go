package address

import "strings"

func IsValidCode(provinceCode string, cityCode string, areaCode string) bool {
	if _, ok := ProvinceCode[provinceCode]; !ok {
		return false
	}

	if _, ok := CityCode[cityCode]; !ok {
		return false
	}

	if _, ok := CountryCode[areaCode]; !ok {
		return false
	}

	pCode := provinceCode[:2]

	if strings.HasPrefix(cityCode, pCode) == false {
		return false
	}

	cCode := cityCode[:4]

	if strings.HasPrefix(areaCode, cCode) == false {
		return false
	}

	return true
}
