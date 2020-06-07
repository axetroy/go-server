// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin_server

import (
	"context"
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/redis"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

func Serve(host string, port string) error {
	redis.Connect()

	defer func() {
		redis.Dispose()
	}()

	database.Connect()

	defer func() {
		database.Dispose()
	}()

	s := &http.Server{
		Addr:           net.JoinHostPort(host, port),
		Handler:        AdminRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 10M
	}

	var wg sync.WaitGroup
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	exit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(exit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-exit
		wg.Add(1)

		config.Common.Exiting = true

		//使用context控制srv.Shutdown的超时时间
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		defer cancel()
		err := s.Shutdown(ctx)
		if err != nil {
			log.Println(err)
		}
		wg.Done()
	}()

	log.Printf("Listen on:  %s\n", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		if err == http.ErrServerClosed {
			log.Println("HTTP 服务已被关闭")
		} else {
			return err
		}
	}

	return nil
}
