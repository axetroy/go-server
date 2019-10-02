// Copyright 2019 Axetroy. All rights reserved. MIT license.
package main

import (
	"fmt"
	App "github.com/axetroy/go-server"
	"github.com/axetroy/go-server/src/config"
	"github.com/axetroy/go-server/src/helper/daemon"
	"github.com/axetroy/go-server/src/server/message_queue_server"
	"github.com/axetroy/go-server/src/util"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Usage = "message queue server controller"
	app.Author = App.Author
	app.Email = App.Email
	app.Version = App.Version
	cli.AppHelpTemplate = App.CliTemplate

	c := make(chan os.Signal)

	signal.Notify(c, os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGTSTP, syscall.SIGUSR1, syscall.SIGUSR2)

	go func() {
		for s := range c {
			switch s {
			case os.Kill, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGTSTP:
				fmt.Println("接收到终止信号, 正在退出进程...")
				config.Common.Exiting = true
				time.AfterFunc(5*time.Second, func() {
					os.Exit(0)
				})
			default:
				fmt.Println("接收到信号:", s)
			}
		}
	}()

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start message queue server",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "daemon, d",
					Usage: "running in daemon mode",
				},
			},
			Action: func(c *cli.Context) error {
				// 判断当其是否是子进程，当父进程return之后，子进程会被系统1号进程接管
				return daemon.Start(message_queue_server.Serve, c.Bool("daemon"))
			},
		},
		{
			Name:  "stop",
			Usage: "stop message queue server",
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
