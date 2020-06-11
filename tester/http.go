package tester

import (
	"github.com/axetroy/go-server/internal/app/admin_server"
	"github.com/axetroy/go-server/internal/app/customer_service"
	"github.com/axetroy/go-server/internal/app/user_server"
	"github.com/axetroy/mocker"
)

var (
	HttpUser            = mocker.New(user_server.UserRouter)                 // 用户接口的模拟器
	HttpAdmin           = mocker.New(admin_server.AdminRouter)               // 管理员接口的模拟器
	HttpCustomerService = mocker.New(customer_service.CustomerServiceRouter) // 客服接口
)
