package address_test

import (
	"github.com/axetroy/go-server/internal/app/user_server/controller/address"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsValidCode(t *testing.T) {
	{
		var codes = [][]string{
			{"123", "123", "123"},          // 不匹配
			{"110000", "130200", "110105"}, // 地区对不上
		}

		for _, arr := range codes {
			provinceCode := arr[0]
			cityCode := arr[1]
			areaCode := arr[2]

			assert.False(t, address.IsValidCode(provinceCode, cityCode, areaCode))
		}
	}

	{
		var codes = [][]string{
			{"110000", "110100", "110101"}, // 北京市 - 北京市 - 东城区
			{"450000", "450100", "450105"}, // 广西 - 南宁市 - 青秀区
		}

		for _, arr := range codes {
			provinceCode := arr[0]
			cityCode := arr[1]
			areaCode := arr[2]

			assert.True(t, address.IsValidCode(provinceCode, cityCode, areaCode))
		}
	}
}
