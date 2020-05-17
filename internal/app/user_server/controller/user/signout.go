// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package user

import (
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/schema"
)

var SignOut = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return schema.Response{}
	})
})
