package token

import (
	"github.com/axetroy/go-server/exception"
)

var (
	ErrInvalidAuth  = exception.NewError("无效的身份认证方式")
	ErrInvalidToken = exception.NewError("无效的身份令牌")
	ErrTokenExpired = exception.NewError("身份令牌已过期")
)
