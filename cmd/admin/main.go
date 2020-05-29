// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package main

import (
	"github.com/axetroy/go-server/internal/app/admin_server"
	"github.com/axetroy/go-server/internal/library/daemon"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Usage = "管理员接口服务器"

	app.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "启动服务",
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
					EnvVars:     []string{"HOST"},
				},
				&cli.StringFlag{
					Name:        "port",
					Usage:       "监听指定端口",
					DefaultText: "80",
					Value:       "80",
					EnvVars:     []string{"PORT"},
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
					return admin_server.Serve(c.String("host"), c.String("port"))
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
		{
			Name:  "migrate",
			Usage: "同步数据库",
			Action: func(context *cli.Context) error {
				return database.Migrate(nil)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
