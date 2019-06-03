package middleware

import (
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/rbac"
	"github.com/axetroy/go-server/src/rbac/accession"
	"github.com/axetroy/go-server/src/schema"
	"github.com/gin-gonic/gin"
	"net/http"
)

// 根据 RBAC 鉴权的中间件
func RequireAccessions(accesions ...accession.Accession) gin.HandlerFunc {
	return func(context *gin.Context) {
		var (
			err error
			uid = context.GetString(ContextUidField) // 这个中间件必须安排在JWT的中间件后面, 所以这里是拿的到 UID 的
			c   *rbac.Controller
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

		if uid == "" {
			err = exception.NoPermission
		}

		if c, err = rbac.New(uid); err != nil {
			return
		}

		if c.Require(accesions) == false {
			err = exception.NoPermission
		}
	}
}
