// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package cache

import (
	"time"
)

type Cache interface {
	GetItem(key string) (interface{}, error)
	SetItem(key string, value interface{}, exp time.Duration) error
	RemoveItem(key string) error
	Clear() error
	Length() (uint, error)
}

var (
//Memory Cache = NewMemoryCache()
//Redis  Cache = NewRedisCache()
)
