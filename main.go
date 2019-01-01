package main

import (
	"github.com/axetroy/go-server/env"
	"github.com/axetroy/go-server/router"
)

func init() {
	if err := env.Load(); err != nil {
		panic(err)
	}
}

func main() {
	if err := router.Router.Run(":8080"); err != nil {
		panic(err)
	}
}
