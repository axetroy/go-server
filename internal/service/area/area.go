// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"sort"
	"strings"
)

var GetArea = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return schema.Response{
			Status:  schema.StatusSuccess,
			Message: "",
			Data:    Maps,
		}
	})
})

var GetDetail = router.Handler(func(c router.Context) {
	fullAreaCode := c.Param("code")

	target, err := LookUp(fullAreaCode)

	c.ResponseFunc(err, func() schema.Response {
		target.Children = nil
		return schema.Response{
			Status:  schema.StatusSuccess,
			Message: "",
			Data:    target,
		}
	})
})

var GetProvinces = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		var result = make([]Location, 0)

		for code, name := range ProvinceMap {
			result = append(result, Location{
				Name: name,
				Code: code,
			})
		}

		sort.SliceStable(result, func(i, j int) bool { return result[i].Code < result[j].Code })

		return schema.Response{
			Status:  schema.StatusSuccess,
			Message: "",
			Data:    result,
		}
	})
})

var GetChildren = router.Handler(func(c router.Context) {
	fullAreaCode := c.Param("code")

	target, err := LookUp(fullAreaCode)

	c.ResponseFunc(err, func() schema.Response {
		var result = make([]Location, 0)

		var resultMap map[string]string

		switch len(target.Code) {
		case 2:
			resultMap = CityMap
			break
		case 4:
			resultMap = AreaMap
			break
		case 6:
			resultMap = StreetMap
			break
		case 9:
			resultMap = nil
			break
		}

		if resultMap != nil {
			for code, name := range resultMap {
				if strings.HasPrefix(code, target.Code) {
					result = append(result, Location{
						Code: code,
						Name: name,
					})
				}
			}
		}

		sort.SliceStable(result, func(i, j int) bool { return result[i].Code < result[j].Code })

		return schema.Response{
			Status:  schema.StatusSuccess,
			Message: "",
			Data:    result,
		}
	})
})
