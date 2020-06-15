// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package main

import (
	"github.com/axetroy/go-server/internal/app/user_server"
	"github.com/axetroy/go-server/pkg/daemon"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "用户接口服务"

	app.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "开启服务",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:        "daemon",
					Usage:       "是否以守护进程运行",
					DefaultText: "false",
					Value:       false,
				},
				&cli.StringFlag{
					Name:        "host",
					Usage:       "监听指定地址",
					DefaultText: "127.0.0.1",
					Value:       "127.0.0.1",
				},
				&cli.StringFlag{
					Name:        "port",
					Usage:       "监听指定端口",
					DefaultText: "80",
					Value:       "80",
				},
				&cli.StringFlag{
					Name:        "domain",
					Usage:       "指定域名",
					DefaultText: "https://example.com",
					Value:       "https://example.com",
					EnvVars:     []string{"DOMAIN"},
				},
			},
			Action: func(c *cli.Context) error {
				// 判断当其是否是子进程，当父进程return之后，子进程会被系统1号进程接管
				return daemon.Start(func() error {
					return user_server.Serve(c.String("host"), c.String("port"))
				}, c.Bool("daemon"))
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
