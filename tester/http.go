package tester

import (
	"github.com/axetroy/go-server/router"
	"github.com/axetroy/mocker"
)

var (
	HttpUser  = mocker.New(router.UserRouter)
	HttpAdmin = mocker.New(router.AdminRouter)
)
