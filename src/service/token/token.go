package token

import (
	"github.com/dgrijalva/jwt-go"
)

var (
	userSecreteKey  = "user"
	adminSecreteKey = "admin"
	Prefix          = "Bearer"
	AuthField       = "Authorization"
)

type Claims struct {
	Uid string `json:"uid"`
	jwt.StandardClaims
}

type ClaimsInternal struct {
	Uid string `json:"uid"` // base64 encode
	jwt.StandardClaims
}
