// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/area"
)

var GetDetail = router.Handler(func(c router.Context) {
	fullAreaCode := c.Param("code")

	target, err := area.LookUp(fullAreaCode)

	c.ResponseFunc(err, func() schema.Response {
		target.Children = nil
		return schema.Response{
			Status:  schema.StatusSuccess,
			Message: "",
			Data:    target,
		}
	})
})
