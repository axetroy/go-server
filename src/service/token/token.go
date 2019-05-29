package token

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"os"
)

var (
	userSecreteKey  = "user"
	adminSecreteKey = "admin"
)

const (
	Prefix    = "Bearer"
	AuthField = "Authorization"
)

type Claims struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

type ClaimsInternal struct {
	Uid string `json:"uid"` // base64 encode
	jwt.StandardClaims
}

func init() {
	userKey := os.Getenv("USER_TOKEN_SECRET_KEY")
	adminKey := os.Getenv("ADMIN_TOKEN_SECRET_KEY")

	if userKey != "" {
		userSecreteKey = userKey
	}

	if adminKey != "" {
		adminSecreteKey = adminKey
	}

	if userSecreteKey == adminSecreteKey {
		panic(errors.New("用户端的 Token 密钥不能和管理员端的相同，存在安全风险"))
	}
}
