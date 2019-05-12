package src

import (
	"fmt"
	"github.com/axetroy/go-server/src/util"
)

func init() {
	if err := util.LoadEnv(); err != nil {
		panic(err)
	}
}

// Server 运行服务器
func ServerUserClient() {
	port := "8080"
	if err := RouterUserClient.Run(":" + port); err != nil {
		panic(err)
	}
	fmt.Println("Listen on port " + port)
}

// Server 运行服务器
func ServerAdminClient() {
	port := "8081"
	if err := RouterAdminClient.Run(":" + port); err != nil {
		panic(err)
	}
	fmt.Println("Listen on port " + port)
}
