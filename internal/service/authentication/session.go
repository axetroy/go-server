// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package authentication

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/service/token"
	nativeRedis "github.com/go-redis/redis"
	"strings"
	"time"
)

type Session struct {
	isAdmin bool
	client  *nativeRedis.Client
}

func NewSession(isAdmin bool) Session {
	var client *nativeRedis.Client

	if isAdmin {
		client = redis.ClientTokenAdmin
	} else {
		client = redis.ClientTokenUser
	}

	return Session{
		isAdmin: isAdmin,
		client:  client,
	}
}

func (c Session) getState() token.State {
	var state token.State
	if c.isAdmin {
		state = token.StateAdmin
	} else {
		state = token.StateUser
	}

	return state
}

func (c Session) Generate(uid string, durations ...time.Duration) (string, error) {
	state := c.getState()

	tokenStr, err := token.Generate(uid, state)

	if err != nil {
		return "", err
	}

	var maxDuration = time.Hour * 24 * 30 // 最长的 token 为一个月
	var duration = time.Hour * 24

	if len(durations) > 0 {
		if durations[0] > 0 {
			duration = durations[0]
		}
	}

	if duration > maxDuration {
		duration = maxDuration
	}

	// 以 token 为 key
	if err := c.client.Set(tokenStr, uid, duration).Err(); err != nil {
		return "", err
	}

	// 以 user_id 为 key
	if err := c.client.Set("id-"+uid+util.RandomString(6), tokenStr, duration).Err(); err != nil {
		return "", err
	}

	return tokenStr, nil
}

func (c Session) Parse(tokenString string) (string, error) {
	tokenString = strings.TrimPrefix(tokenString, token.Prefix+" ")

	uid, err := c.client.Get(tokenString).Result()

	if err != nil {
		return "", exception.InvalidToken
	}

	return uid, nil
}

func (c Session) Remove(tokenString string) error {
	return c.client.Del(tokenString).Err()
}
