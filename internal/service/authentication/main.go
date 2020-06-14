// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package authentication

import "time"

// 该模块定义了 如何生成/解析 token

type Authentication interface {
	Generate(uid string, duration ...time.Duration) (string, error) // 生成身份认证的 token
	Parse(token string) (string, error)                             // 解析 token，并且返回 uuid
	Remove(token string) error                                      // 移除某个 token
}

var JwtUser Authentication      // JWT 的认证方式
var JwtAdmin Authentication     // JWT 的认证方式
var SessionUser Authentication  // JWT 的认证方式
var SessionAdmin Authentication // JWT 的认证方式

func init() {
	JwtUser = Jwt{}
	JwtAdmin = Jwt{IsAdmin: true}
	SessionUser = Session{}
	SessionAdmin = Session{IsAdmin: true}
}
