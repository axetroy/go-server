// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package address

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/jinzhu/gorm"
	"regexp"
)

func AreaList() (res schema.Response) {
	var (
		err  error
		data schema.Area
		tx   *gorm.DB
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	data = schema.Area{
		Province: ProvinceCode,
		City:     CityCode,
		Area:     CountryCode,
	}

	return
}

var AreaListRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return AreaList()
	})
})

type AreaStruct struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Area struct {
	Province AreaStruct `json:"province"`
	City     AreaStruct `json:"city"`
	Country  AreaStruct `json:"country"`
	Addr     string     `json:"addr"`
}

var (
	areaCodeReg = regexp.MustCompile("^(\\d{2})(\\d{2})(\\d{2})$")
)

// 查找地区码相关信息
func FindAddress(areaCode string) (*Area, error) {
	countryName, ok := CountryCode[areaCode]

	if !ok {
		return nil, errors.New(fmt.Sprintf("Invalid code: %s", areaCode))
	}

	matcher := areaCodeReg.FindAllStringSubmatch(areaCode, 1)

	s := matcher[0]

	pCode := s[1] + "0000"      // 省份代码
	cCode := s[1] + s[2] + "00" // 城市代码

	provinceName := ProvinceCode[pCode]
	cityName := CityCode[cCode]

	oCityName := cityName

	if cityName == provinceName {
		cityName = ""
	}

	addr := provinceName + cityName + countryName

	area := Area{
		Province: AreaStruct{Code: pCode, Name: provinceName},
		City:     AreaStruct{Code: cCode, Name: oCityName},
		Country:  AreaStruct{Code: areaCode, Name: countryName},
		Addr:     addr,
	}

	return &area, nil
}

var FindAddressRouter = router.Handler(func(c router.Context) {
	var (
		err  error
		res  = schema.Response{}
		area *Area
	)

	code := c.Param("area_code")

	area, err = FindAddress(code)

	res.Data = area

	c.ResponseFunc(err, func() schema.Response {
		return res
	})
})
