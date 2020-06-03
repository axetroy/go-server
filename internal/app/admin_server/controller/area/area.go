// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/area"
)

type query struct {
	Simple  bool `json:"simple" url:"simple"`   // 简化字段输出
	Flatten bool `json:"flatten" url:"flatten"` // 扁平化输出
}

var GetArea = router.Handler(func(c router.Context) {
	var query query

	c.ResponseFunc(c.ShouldBindQuery(&query), func() schema.Response {
		if query.Flatten {
			return schema.Response{
				Status:  schema.StatusSuccess,
				Message: "",
				Data:    area.MapsFlatten,
			}
		} else if query.Simple {
			return schema.Response{
				Status:  schema.StatusSuccess,
				Message: "",
				Data:    area.MapsSimplified,
			}
		} else {
			return schema.Response{
				Status:  schema.StatusSuccess,
				Message: "",
				Data:    area.Maps,
			}
		}
	})
})
