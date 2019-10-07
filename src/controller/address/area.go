// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/schema"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
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

		helper.Response(&res, data, err)
	}()

	data = schema.Area{
		Province: ProvinceCode,
		City:     CityCode,
		Area:     CountryCode,
	}

	return
}

func AreaListRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	res = AreaList()
}

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

func FindAddressRouter(c *gin.Context) {
	var (
		err  error
		res  = schema.Response{}
		area *Area
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	code := c.Param("area_code")

	area, err = FindAddress(code)

	res.Data = area
}
