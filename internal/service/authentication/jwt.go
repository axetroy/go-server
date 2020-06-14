// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package authentication

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/service/token"
	"time"
)

type Jwt struct {
	IsAdmin bool
}

func (c Jwt) getState() token.State {
	var state token.State
	if c.IsAdmin {
		state = token.StateAdmin
	} else {
		state = token.StateUser
	}

	return state
}

func (c Jwt) Generate(uid string, duration ...time.Duration) (string, error) {
	state := c.getState()

	return token.Generate(uid, state, duration...)
}

func (c Jwt) Parse(tokenString string) (string, error) {
	state := c.getState()

	if claims, err := token.Parse(tokenString, state); err != nil {
		return "", exception.InvalidToken
	} else {
		return claims.Uid, nil
	}
}

func (c Jwt) Remove(_ string) error {
	return nil
}
