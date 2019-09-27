// Copyright 2019 Axetroy. All rights reserved. MIT license.
package main

import (
	"fmt"
	"github.com/axetroy/go-fs"
	"github.com/axetroy/go-server/src"
	"github.com/urfave/cli"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

var (
	pidFileName = "go-server.admin.pid"
)

func main() {
	app := cli.NewApp()
	app.Usage = "admin server controller"

	app.Commands = []cli.Command{
		{
			Name:  "start",
			Usage: "start admin server",
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "daemon, d",
					Usage: "running in daemon mode",
				},
			},
			Action: func(c *cli.Context) error {
				// 判断当其是否是子进程，当父进程return之后，子进程会被系统1号进程接管
				daemon := c.Bool("daemon")

				if daemon && os.Getppid() != 1 {
					// 将命令行参数中执行文件路径转换成可用路径
					filePath, _ := filepath.Abs(os.Args[0])
					cmd := exec.Command(filePath, os.Args[1:]...)
					// 将其他命令传入生成出的进程
					//cmd.Stdin = os.Stdin // 给新进程设置文件描述符，可以重定向到文件中
					//cmd.Stdout = os.Stdout
					//cmd.Stderr = os.Stderr
					_ = cmd.Start() // 开始执行新进程，不等待新进程退出
					return nil
				} else {
					// 讲 pid 写入到当前文件夹下
					if err := fs.WriteFile(pidFileName, []byte(fmt.Sprintf("%d", os.Getpid()))); err != nil {
						return err
					}
					src.ServerAdminClient()
					return nil
				}
			},
		},
		{
			Name:  "stop",
			Usage: "stop admin server",
			Action: func(c *cli.Context) error {
				if !fs.PathExists(pidFileName) {
					return nil
				}

				b, err := fs.ReadFile(pidFileName)

				if err != nil {
					return nil
				}

				pidStr := string(b)

				pid, err := strconv.Atoi(pidStr)

				if err != nil {
					return err
				}

				ps, err := os.FindProcess(pid)

				if err != nil {
					return err
				}

				if err := ps.Kill(); err != nil {
					return err
				}

				fmt.Printf("process %s have been kill.\n", pidStr)

				_ = fs.Remove(pidFileName)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
