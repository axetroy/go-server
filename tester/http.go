package tester

import (
	"github.com/axetroy/go-server/internal/app/admin_server"
	"github.com/axetroy/go-server/internal/app/user_server"
	"github.com/axetroy/mocker"
)

var (
	//HttpUser 用户接口的模拟器
	HttpUser = mocker.New(user_server.UserRouter)
	//HttpAdmin 管理员接口的模拟器
	HttpAdmin = mocker.New(admin_server.AdminRouter)
)
