package middleware

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/axetroy/server/response"
	"gitlab.com/axetroy/server/token"
	"net/http"
)

// 普通用户的token验证
func Authenticate() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")

		if claims, err := token.Parse(tokenString); err != nil {
			context.JSON(http.StatusOK, response.Response{
				Message: err.Error(),
				Data:    nil,
			})

			context.Abort()

			return
		} else {
			context.Set("uid", claims.Uid)
		}
	}
}

// TODO: 管理员的token验证
func AuthenticateAdmin() gin.HandlerFunc {
	return func(context *gin.Context) {
		tokenString := context.GetHeader("Authorization")

		if claims, err := token.Parse(tokenString); err != nil {
			context.JSON(http.StatusOK, response.Response{
				Message: err.Error(),
				Data:    nil,
			})

			context.Abort()

			return
		} else {
			context.Set("uid", claims.Uid)
		}
	}
}