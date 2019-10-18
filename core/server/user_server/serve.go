// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user_server

import (
	"context"
	"crypto/tls"
	"github.com/axetroy/go-server/core/config"
	"github.com/axetroy/go-server/core/service/database"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func Serve() error {
	port := config.User.Port

	s := &http.Server{
		Addr:           ":" + port,
		Handler:        UserRouter,
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20, // 10M
	}

	log.Printf("Listen on:  %s\n", s.Addr)

	go func() {
		if config.User.TLS != nil {
			TLSConfig := &tls.Config{
				MinVersion:               tls.VersionTLS11,
				CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
				PreferServerCipherSuites: true,
				CipherSuites: []uint16{
					tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
					tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
					tls.TLS_RSA_WITH_AES_256_CBC_SHA,
				},
			}

			TLSProto := make(map[string]func(*http.Server, *tls.Conn, http.Handler), 0)

			s.TLSConfig = TLSConfig
			s.TLSNextProto = TLSProto

			if err := s.ListenAndServeTLS(config.User.TLS.Cert, config.User.TLS.Key); err != nil {
				log.Println(err)
			}
		} else {
			if err := s.ListenAndServe(); err != nil {
				log.Println(err)
			}
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTSTP)
	<-quit

	config.Common.Exiting = true

	log.Println("Shutdown Server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}

	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Println("Timeout of 5 seconds.")
	}

	_ = database.Db.Close()

	log.Println("Server exiting")

	return nil
}
