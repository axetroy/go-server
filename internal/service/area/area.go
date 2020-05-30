// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
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

var GetChildren = router.Handler(func(c router.Context) {
	fullAreaCode := c.Param("code")

	target, err := LookUp(fullAreaCode)

	c.ResponseFunc(err, func() schema.Response {
		var result = make([]Location, 0)
		for _, location := range target.Children {
			location.Children = nil
			result = append(result, location)
		}
		return schema.Response{
			Status:  schema.StatusSuccess,
			Message: "",
			Data:    result,
		}
	})
})
