package tester

import (
	"github.com/axetroy/go-server/src/server/admin_server"
	"github.com/axetroy/go-server/src/server/user_server"
	"github.com/axetroy/mocker"
)

var (
	HttpUser  = mocker.New(user_server.UserRouter)
	HttpAdmin = mocker.New(admin_server.AdminRouter)
)
