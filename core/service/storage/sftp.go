// Copyright 2019 Axetroy. All rights reserved. MIT license.
package storage

func NewSFTPStorage() *SFTPStorage {
	c := &SFTPStorage{}

	return c
}

type SFTPStorage struct {
}

func (c *SFTPStorage) Store(file []byte, filename string) (*string, error) {
	return nil, nil
}
