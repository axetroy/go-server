package main

import (
	"github.com/joho/godotenv"
	"gitlab.com/axetroy/server/router"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	if err := router.Router.Run(":8080"); err != nil {
		panic(err)
	}
}
