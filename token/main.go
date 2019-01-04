package token

import (
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/utils"
	"github.com/dgrijalva/jwt-go"
	"strings"
	"time"
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

// generate jwt token
func Generate(userId string, isAdmin bool) (tokenString string, err error) {
	var (
		issuer string
		key    string
	)

	if isAdmin {
		issuer = "admin"
		key = adminSecreteKey
	} else {
		issuer = "user"
		key = userSecreteKey
	}

	// 生成token
	c := ClaimsInternal{
		utils.Base64Encode(userId),
		jwt.StandardClaims{
			Audience:  userId,
			Id:        userId,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(6)).Unix(),
			Issuer:    issuer,
			IssuedAt:  time.Now().Unix(),
			NotBefore: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err = token.SignedString([]byte(key))

	return
}

// parse jwt token
func Parse(tokenString string, isAdmin bool) (claims Claims, err error) {
	var (
		token *jwt.Token
		key   string
	)

	if isAdmin {
		key = adminSecreteKey
	} else {
		key = userSecreteKey
	}

	if strings.HasPrefix(tokenString, Prefix+" ") == false {
		err = exception.InvalidAuth
		return
	}

	tokenString = strings.Replace(tokenString, Prefix+" ", "", 1)

	if tokenString == "" {
		err = exception.InvalidToken
		return
	}

	c := ClaimsInternal{}

	if token, err = jwt.ParseWithClaims(tokenString, &c, func(token *jwt.Token) (interface{}, error) {
		return []byte(key), nil
	}); err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			err = exception.TokenExpired
		}
		err = exception.InvalidToken
		return
	}

	if token != nil && token.Valid {
		var (
			uid string
		)

		if uid, err = utils.Base64Decode(c.Uid); err != nil {
			return
		}

		claims.Uid = uid
		claims.Audience = c.Audience
		claims.Id = c.Id
		claims.NotBefore = c.NotBefore
		claims.ExpiresAt = c.ExpiresAt
		claims.Issuer = c.Issuer
		claims.IssuedAt = c.IssuedAt
		claims.Subject = c.Subject

		return
	} else if ve, ok := err.(*jwt.ValidationError); ok {
		if ve.Errors&jwt.ValidationErrorMalformed != 0 {
			err = exception.InvalidToken
			return
		} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
			err = exception.TokenExpired
			return
		} else {
			err = exception.InvalidToken
			return
		}
	} else {
		err = exception.InvalidToken
		return
	}
}
