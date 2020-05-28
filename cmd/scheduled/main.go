// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package main

import (
	"github.com/axetroy/go-server/cmd/scheduled/job"
	"github.com/axetroy/go-server/internal/library/daemon"
	"github.com/jasonlvhit/gocron"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func runJobs() error {
	// 每天凌晨 3 点检查 login_log 表，并且进行切割数据
	// 选择半夜主要是因为怕影响性能，在用户最少的情况下执行
	if err := gocron.Every(1).Day().At("03:00:01").Do(job.SplitLoginLog); err != nil {
		return err
	}

	// 启动定时任务
	<-gocron.Start()

	return nil
}

func main() {
	app := cli.NewApp()
	app.Usage = "定时任务"

	app.Commands = []*cli.Command{
		{
			Name:  "start",
			Usage: "启动定时任务",
			Flags: []cli.Flag{
				&cli.BoolFlag{
					Name:  "daemon, d",
					Usage: "是否以守护进程运行",
				},
			},
			Action: func(c *cli.Context) error {
				// 判断当其是否是子进程，当父进程return之后，子进程会被系统1号进程接管
				return daemon.Start(runJobs, c.Bool("daemon"))
			},
		},
		{
			Name:  "stop",
			Usage: "停止定时任务",
			Action: func(c *cli.Context) error {
				return daemon.Stop()
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
