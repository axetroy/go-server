package middleware

import (
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/gin-gonic/gin"
	"net/http"
)

var (
	ContextUidField = "uid"
)

// Token 验证中间件
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
			// 把 UID 挂载到上下文中国呢
			context.Set(ContextUidField, claims.Uid)
		}
	}
}
