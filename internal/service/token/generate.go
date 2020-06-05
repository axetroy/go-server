// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package token

import (
	"github.com/axetroy/go-server/internal/library/config"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type State string

const (
	StateUser  State = "user"
	StateAdmin State = "admin"
)

// generate jwt token
func Generate(userId string, state State, d ...time.Duration) (tokenString string, err error) {
	var (
		issuer string
		key    string
	)

	var duration time.Duration

	if len(d) > 0 {
		duration = d[0]
		if d[0] == 0 {
			duration = time.Hour * time.Duration(6)
		}
	} else {
		duration = time.Hour * time.Duration(6)
	}

	// Token 有效期最高 30 天
	if duration > time.Hour*24*30 {
		duration = time.Hour * 24 * 30
	}

	switch state {
	case StateAdmin:
		issuer = "admin"
		key = config.Jwt.Secret
		break
	case StateUser:
		issuer = "user"
		key = config.Jwt.Secret
		break
	}

	// 生成token
	c := ClaimsInternal{
		util.Base64Encode(userId),
		jwt.StandardClaims{
			Audience:  userId,
			Id:        userId,
			ExpiresAt: time.Now().Add(duration).Unix(),
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err = token.SignedString([]byte(key))

	return
}
