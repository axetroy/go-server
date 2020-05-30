// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package area

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/area"
)

var GetArea = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return schema.Response{
			Status:  schema.StatusSuccess,
			Message: "",
			Data:    area.Maps,
		}
	})
})
