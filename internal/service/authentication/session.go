// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package authentication

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/service/token"
	"time"
)

type Session struct {
	IsAdmin bool
}

func (c Session) getState() token.State {
	var state token.State
	if c.IsAdmin {
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

	var duration = time.Hour * 24

	if len(durations) > 0 {
		duration = durations[0]
	}

	if state == token.StateUser {
		if err := redis.ClientTokenUser.Set(tokenStr, uid, duration).Err(); err != nil {
			return "", err
		}
	} else {
		if err := redis.ClientTokenAdmin.Set(tokenStr, uid, duration).Err(); err != nil {
			return "", err
		}
	}

	return tokenStr, nil
}

func (c Session) Parse(tokenString string) (string, error) {
	state := c.getState()

	if state == token.StateUser {
		uid, err := redis.ClientTokenUser.Get(tokenString).Result()

		if err != nil {
			return "", exception.InvalidToken
		}

		return uid, nil
	} else {
		uid, err := redis.ClientTokenAdmin.Get(tokenString).Result()

		if err != nil {
			return "", exception.InvalidToken
		}

		return uid, nil
	}
}

func (c Session) Remove(tokenString string) error {
	state := c.getState()

	if state == token.StateUser {
		return redis.ClientTokenUser.Del(tokenString).Err()
	} else {
		return redis.ClientTokenAdmin.Del(tokenString).Err()
	}
}
