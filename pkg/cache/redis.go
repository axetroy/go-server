// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package cache

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

type Redis struct {
	c *redis.Client
}

func NewRedisCache(addr string, password string) Redis {
	return Redis{
		c: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       13,
		}),
	}
}

func (m Redis) GetItem(key string) (interface{}, error) {
	val, err := m.c.WithTimeout(time.Minute).Get(context.Background(), key).Result()

	if err != nil {
		return nil, err
	}

	return val, nil
}

func (m Redis) SetItem(key string, value interface{}, exp time.Duration) error {
	err := m.c.WithTimeout(time.Minute).Set(context.Background(), key, value, exp).Err()

	if err != nil {
		return err
	}

	return nil
}

func (m Redis) RemoveItem(key string) error {
	if err := m.c.WithTimeout(time.Minute).Del(context.Background(), key).Err(); err != nil {
		return err
	}

	return nil
}

func (m Redis) Clear() error {
	// TODO
	return nil
}

func (m Redis) Length() (uint, error) {
	// TODO
	return uint(0), nil
}
