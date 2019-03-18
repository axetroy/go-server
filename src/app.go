package src

import (
	"fmt"
	"github.com/axetroy/go-server/src/controller/admin"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
)

func init() {
	if err := util.LoadEnv(); err != nil {
		panic(err)
	}
}

// Init 初始化一些动作
func Init() {
	// 确保超级管理员存在
	r := admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "admin",
		Password: "admin",
		Name:     "admin",
	}, true)

	// 如果抛出错误，并且不是超级管理员已存在的话
	if r.Status != schema.StatusSuccess && r.Message != exception.AdminExist.Error() {
		panic(r.Message)
	}
}

// Server 运行服务器
func Server() {
	Init()
	port := "8080"
	if err := Router.Run(":" + port); err != nil {
		panic(err)
	}
	fmt.Println("Listen on port " + port)
}
