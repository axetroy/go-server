package tester

import (
	"github.com/axetroy/go-server/src"
	"github.com/axetroy/mocker"
)

var (
	HttpUser  = mocker.New(src.UserRouter)
	HttpAdmin = mocker.New(src.AdminRouter)
)
