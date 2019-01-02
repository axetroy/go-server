package main

import (
	"github.com/axetroy/go-server/controller/admin"
	"github.com/axetroy/go-server/env"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/router"
)

func init() {
	if err := env.Load(); err != nil {
		panic(err)
	}
}

func main() {
	// 确保超级管理员存在
	r := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "admin",
		Password: "admin",
		Name:     "admin",
	}, true)

	// 如果抛出错误，并且不是超级管理员已存在的话
	if r.Status != response.StatusSuccess && r.Message != exception.AdminExist.Error() {
		panic(r.Message)
	}

	if err := router.Router.Run(":8080"); err != nil {
		panic(err)
	}
}
