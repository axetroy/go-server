package token

import (
	"github.com/axetroy/go-server/common_error"
)

var (
	ErrInvalidAuth  = common_error.NewError("无效的身份认证方式")
	ErrInvalidToken = common_error.NewError("无效的身份令牌")
	ErrTokenExpired = common_error.NewError("身份令牌已过期")
)
