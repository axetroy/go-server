package src

import (
	"fmt"
	"github.com/axetroy/go-server/src/config"
	"github.com/axetroy/go-server/src/service/dotenv"
	"net/http"
	"time"
)

func init() {
	if err := dotenv.Load(); err != nil {
		panic(err)
	}
}

// Server 运行服务器
func ServerUserClient() {
	port := config.User.Port

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        UserRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1024 * 1024 * 20, // 20M
	}
	fmt.Printf("用户端 HTTP 监听:  %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}

// Server 运行服务器
func ServerAdminClient() {
	port := config.Admin.Port

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        AdminRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1024 * 1024 * 20, // 20M
	}
	fmt.Printf("管理员端 HTTP 监听:  %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
}
