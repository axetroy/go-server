// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package report

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
)

var GetTypesRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return schema.Response{
			Message: "",
			Status:  schema.StatusSuccess,
			Data:    model.ReportTypes,
		}
	})
})
