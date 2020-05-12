// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package main

import (
	App "github.com/axetroy/go-server"
	message_queue_server2 "github.com/axetroy/go-server/internal/app/message_queue_server"
	"github.com/axetroy/go-server/internal/library/daemon"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/urfave/cli"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "message queue server"
	app.Author = App.Author
	app.Email = App.Email
	app.Version = App.Version
	cli.AppHelpTemplate = App.CliTemplate

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start server",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "daemon, d",
					Usage: "running in daemon mode",
				},
			},
			Action: func(c *cli.Context) error {
				// 判断当其是否是子进程，当父进程return之后，子进程会被系统1号进程接管
				return daemon.Start(message_queue_server2.Serve, c.Bool("daemon"))
			},
		},
		{
			Name:  "stop",
			Usage: "stop message queue",
			Action: func(c *cli.Context) error {
				return daemon.Stop()
			},
		},
		{
			Name:  "env",
			Usage: "print runtime environment",
			Action: func(c *cli.Context) error {
				util.PrintEnv()
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
