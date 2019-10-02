// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin_server

import (
	"fmt"
	"github.com/axetroy/go-server/src/config"
	"net/http"
	"time"
)

func Serve() error {
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
		return err
	}

	return nil
}
