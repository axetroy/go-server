// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package finance

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
)

var GetHistory = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return schema.Response{}
	})
})
