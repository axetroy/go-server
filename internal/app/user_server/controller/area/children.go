// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/area"
	"sort"
	"strings"
)

var GetChildren = router.Handler(func(c router.Context) {
	fullAreaCode := c.Param("code")

	target, err := area.LookUp(fullAreaCode)

	c.ResponseFunc(err, func() schema.Response {
		var result = make([]area.Location, 0)

		var resultMap map[string]string

		switch len(target.Code) {
		case 2:
			resultMap = area.CityMap
		case 4:
			resultMap = area.AreaMap
		case 6:
			resultMap = area.StreetMap
		case 9:
			resultMap = nil
		}

		for code, name := range resultMap {
			if strings.HasPrefix(code, target.Code) {
				result = append(result, area.Location{
					Code: code,
					Name: name,
				})
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
