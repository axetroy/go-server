// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package storage

import (
	"fmt"
	"github.com/axetroy/go-server/internal/library/config"
	"log"
)

type Storage interface {
	Store(file []byte, filename string) (*string, error) // 存储文件
}

type provider string

var (
	client        *Storage           // 文件存储的客户端
	providerLocal provider = "local" // 本地存储
	providerSFTP  provider = "sftp"  // SFTP 协议储存到远端
)

func init() {
	switch provider(config.Telephone.Provider) {
	case providerLocal:
		initClient(NewLocalStorage())
	case providerSFTP:
		initClient(NewSFTPStorage())
	default:
		log.Fatal(fmt.Sprintf(`Invalid storage provider "%s"`, config.Telephone.Provider))
	}
}

func initClient(s Storage) {
	client = &s
}

func GetClient() Storage {
	return *client
}
