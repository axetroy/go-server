package middleware

import (
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
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
				context.JSON(http.StatusOK, schema.Response{
					Message: err.Error(),
					Data:    nil,
				})
				context.Abort()
			}
		}()

		if s, isExist := context.GetQuery(util.AuthField); isExist == true {
			tokenString = s
			return
		} else {
			tokenString = context.GetHeader(util.AuthField)

			if len(tokenString) == 0 {
				if s, er := context.Cookie(util.AuthField); er != nil {
					err = er
					return
				} else {
					tokenString = s
				}
			}
		}

		if claims, er := util.ParseToken(tokenString, isAdmin); er != nil {
			err = er
			return
		} else {
			context.Set("uid", claims.Uid)
		}
	}
}
