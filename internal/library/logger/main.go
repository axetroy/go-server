package logger

import (
	"github.com/axetroy/go-fs"
	log "github.com/sirupsen/logrus"
	"os"
	"path"
)

var (
	Info  = log.Info
	Infof = log.Infof
)

func init() {
	// 设置日志格式为json格式
	log.SetFormatter(&log.JSONFormatter{})

	logBaseDir := ""

	cwd, err := os.Getwd()

	if err != nil {
		exPath, err := os.Executable()

		if err != nil {
			panic(err)
		}

		logBaseDir = exPath
	} else {
		logBaseDir = cwd
	}

	logsPath := path.Join(logBaseDir, "logs")
	logFilePath := path.Join(logsPath, "main.log")

	if err := fs.EnsureDir(logsPath); err != nil {
		panic(err)
	}

	if err := fs.EnsureFile(logFilePath); err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm)

	if err != nil {
		panic(err)
	}

	// 设置将日志输出到标准输出（默认的输出为stderr,标准错误）
	// 日志消息输出可以是任意的io.writer类型
	log.SetOutput(logFile)

	// 设置日志级别为warn以上
	log.SetLevel(log.InfoLevel)
}
