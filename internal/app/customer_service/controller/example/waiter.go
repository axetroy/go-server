// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package example

import (
	_ "github.com/axetroy/go-server/internal/app/customer_service/views_pkged"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/markbates/pkger"
	"io"
)

var WaiterRouter = router.Handler(func(c router.Context) {
	f, err := pkger.Open("/internal/app/customer_service/views/waiter.html")

	if err != nil {
		c.Writer().Header().Del("Content-Security-Policy")
		_, _ = c.Writer().Write([]byte(err.Error()))
		return
	}

	_, _ = io.Copy(c.Writer(), f)
})
