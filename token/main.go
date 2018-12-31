package token

import (
	"github.com/axetroy/redpack/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/axetroy/go-server/exception"
	"strconv"
	"strings"
	"time"
)

var (
	SecreteKey = "hello"
	Prefix     = "Bearer"
)

type Claims struct {
	Uid int64 `json:"uid"`
	jwt.StandardClaims
}

type ClaimsInternal struct {
	Uid string `json:"uid"` // base64 encode
	jwt.StandardClaims
}

// generate jwt token
func Generate(userId int64) (tokenString string, err error) {
	// 生成token
	var idStr = strconv.FormatInt(userId, 10)
	c := ClaimsInternal{
		utils.Base64Encode(idStr),
		jwt.StandardClaims{
			Audience:  idStr,
			Id:        idStr,
			ExpiresAt: time.Now().Add(time.Hour * time.Duration(6)).Unix(),
			Issuer:    "test",
			IssuedAt:  time.Now().Unix(),
			Subject:   "test",
			NotBefore: time.Now().Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)

	tokenString, err = token.SignedString([]byte(SecreteKey))

	return
}

// parse jwt token
func Parse(tokenString string) (claims Claims, err error) {
	var (
		token *jwt.Token
	)

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
		return []byte(SecreteKey), nil
	}); err != nil {
		if strings.HasPrefix(err.Error(), "token is expired by") {
			err = exception.TokenExpired
		}
		return
	}

	if token != nil && token.Valid {

		var (
			uidStr string
			uid    int64
		)

		if uidStr, err = utils.Base64Decode(c.Uid); err != nil {
			return
		}

		if uid, err = strconv.ParseInt(uidStr, 10, 64); err != nil {
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
