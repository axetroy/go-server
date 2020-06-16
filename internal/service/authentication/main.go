// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package authentication

import (
	"os"
	"time"
)

// 该模块定义了 如何生成/解析 token

type Authentication interface {
	Generate(uid string, duration ...time.Duration) (string, error) // 生成身份认证的 token
	Parse(token string) (string, error)                             // 解析 token，并且返回 uuid
	Remove(token string) error                                      // 移除某个 token
}

var jwtUser Authentication = NewJwt(false)         // JWT 的认证方式
var jwtAdmin Authentication = NewJwt(true)         // JWT 的认证方式
var sessionUser Authentication = NewSession(false) // session 的认证方式
var sessionAdmin Authentication = NewSession(true) // session 的认证方式

var Gateway func(isAdmin bool) Authentication // 网关，所有的认证方式都通过这里

func init() {
	// 通过 JWT 环境变量
	if len(os.Getenv("JWT")) > 0 {
		Gateway = func(isAdmin bool) Authentication {
			if isAdmin {
				return jwtAdmin
			} else {
				return jwtUser
			}
		}
	} else {
		Gateway = func(isAdmin bool) Authentication {
			if isAdmin {
				return sessionAdmin
			} else {
				return sessionUser
			}
		}
	}

}
