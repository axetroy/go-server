// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/area"
	"sort"
)

var GetProvinces = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		var result = make([]area.Location, 0)

		for code, name := range area.ProvinceMap {
			result = append(result, area.Location{
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
