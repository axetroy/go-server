// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package cache

import (
	memCache "github.com/patrickmn/go-cache"
	"time"
)

type Mem struct {
	c *memCache.Cache
}

func NewMemoryCache() Mem {
	return Mem{
		c: memCache.New(time.Minute*10, time.Minute),
	}
}

func (m Mem) GetItem(key string) (interface{}, error) {
	if val, ok := m.c.Get(key); !ok {
		return nil, nil
	} else {
		return val, nil
	}
}

func (m Mem) SetItem(key string, value interface{}, exp time.Duration) error {
	m.c.Set(key, value, exp)

	return nil
}

func (m Mem) RemoveItem(key string) error {
	m.c.Delete(key)

	return nil
}

func (m Mem) Clear() error {
	items := m.c.Items()

	for key := range items {
		m.c.Delete(key)
	}

	return nil
}

func (m Mem) Length() (uint, error) {
	length := len(m.c.Items())

	return uint(length), nil
}
