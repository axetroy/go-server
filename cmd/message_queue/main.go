// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package main

import (
	message_queue_server2 "github.com/axetroy/go-server/internal/app/message_queue_server"
	"github.com/axetroy/go-server/pkg/daemon"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "消息队列服务器"

	app.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "开启服务",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "daemon, d",
					Usage: "是否以守护进程运行",
				},
			},
			Action: func(c *cli.Context) error {
				// 判断当其是否是子进程，当父进程return之后，子进程会被系统1号进程接管
				return daemon.Start(message_queue_server2.Serve, c.Bool("daemon"))
			},
		},
		{
			Name:  "stop",
			Usage: "停止服务",
			Action: func(c *cli.Context) error {
				return daemon.Stop()
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
