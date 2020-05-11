package tester

import (
	admin_server2 "github.com/axetroy/go-server/internal/app/admin_server"
	user_server2 "github.com/axetroy/go-server/internal/app/user_server"
	"github.com/axetroy/mocker"
)

var (
	HttpUser  = mocker.New(user_server2.UserRouter)
	HttpAdmin = mocker.New(admin_server2.AdminRouter)
)
