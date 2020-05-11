// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package storage

func NewLocalStorage() *LocalStorage {
	c := &LocalStorage{}

	return c
}

type LocalStorage struct {
	RootPath string `json:"root_path"` // 定义存储的根目录
}

func (c *LocalStorage) Store(file []byte, filename string) (*string, error) {
	return nil, nil
}
