package middleware

import (
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 普通用户的token验证
func Authenticate(isAdmin bool) gin.HandlerFunc {
	return func(context *gin.Context) {
		var (
			err         error
			tokenString string
		)
		defer func() {
			if err != nil {
				context.JSON(http.StatusOK, response.Response{
					Message: err.Error(),
					Data:    nil,
				})
				context.Abort()
			}
		}()

		if s, isExist := context.GetQuery(token.AuthField); isExist == true {
			tokenString = s
			return
		} else {
			tokenString = context.GetHeader(token.AuthField)

			if len(tokenString) == 0 {
				if s, er := context.Cookie(token.AuthField); er != nil {
					err = er
					return
				} else {
					tokenString = s
				}
			}
		}

		if claims, er := token.Parse(tokenString, isAdmin); er != nil {
			err = er
			return
		} else {
			context.Set("uid", claims.Uid)
		}
	}
}
