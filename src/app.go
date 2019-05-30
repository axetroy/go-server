package src

import (
	"fmt"
	"github.com/axetroy/go-server/src/service/dotenv"
	"net/http"
	"os"
	"time"
)

func init() {
	if err := dotenv.Load(); err != nil {
		panic(err)
	}
}

// Server 运行服务器
func ServerUserClient() {
	port := "8080"

	if p := os.Getenv("USER_HTTP_PORT"); p != "" {
		port = p
	}

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
	port := "8081"

	if p := os.Getenv("ADMIN_HTTP_PORT"); p != "" {
		port = p
	}

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
