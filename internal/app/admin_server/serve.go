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
	"syscall"
	"time"
)

func Serve(host string, port string) error {
	redis.Connect()
	database.Connect()

	s := &http.Server{
		Addr:           net.JoinHostPort(host, port),
		Handler:        AdminRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 10M
	}

	log.Printf("Listen on:  %s\n", s.Addr)

	go func() {
		if err := s.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-quit

	config.Common.Exiting = true

	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	<-ctx.Done()
	log.Println("Timeout of 5 seconds.")

	redis.Dispose()
	database.Dispose()

	log.Println("Server exiting")

	return nil
}
