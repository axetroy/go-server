// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package main

import (
	"github.com/axetroy/go-server/internal/app/user_server"
	"github.com/axetroy/go-server/internal/library/daemon"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "user server"

	app.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "start server",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "daemon, d",
					Usage: "running in daemon mode",
				},
			},
			Action: func(c *cli.Context) error {
				// 判断当其是否是子进程，当父进程return之后，子进程会被系统1号进程接管
				return daemon.Start(user_server.Serve, c.Bool("daemon"))
			},
		},
		{
			Name:  "stop",
			Usage: "stop server",
			Action: func(c *cli.Context) error {
				return daemon.Stop()
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
