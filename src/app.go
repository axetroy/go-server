package src

import (
	"fmt"
	"github.com/axetroy/go-server/src/util"
	"net/http"
	"time"
)

func init() {
	if err := util.LoadEnv(); err != nil {
		panic(err)
	}
}

// Server 运行服务器
func ServerUserClient() {
	port := "8080"
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        UserRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1024 * 1024 * 20, // 20M
	}
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
	fmt.Println("Listen on port " + port)
}

// Server 运行服务器
func ServerAdminClient() {
	port := "8081"
	s := &http.Server{
		Addr:           ":" + port,
		Handler:        AdminRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1024 * 1024 * 20, // 20M
	}
	if err := s.ListenAndServe(); err != nil {
		panic(err)
	}
	fmt.Println("Listen on port " + port)
}
