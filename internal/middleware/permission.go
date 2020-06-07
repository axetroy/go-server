package middleware

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/rbac"
	"github.com/axetroy/go-server/internal/rbac/accession"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/kataras/iris/v12"
)

// 根据 RBAC 鉴权的中间件
func Permission(accessions ...accession.Accession) iris.Handler {
	return func(c iris.Context) {
		var (
			err error
			uid = c.Values().GetString("uid") // 这个中间件必须安排在JWT的中间件后面, 所以这里是拿的到 UID 的
			cc  *rbac.Controller
		)

		defer func() {
			if err != nil {
				_, _ = c.JSON(schema.Response{
					Message: err.Error(),
					Data:    nil,
				})
				return
			}

			c.Next()
		}()

		if uid == "" {
			err = exception.NoPermission
		}

		if cc, err = rbac.New(uid); err != nil {
			return
		}

		if !cc.Require(accessions) {
			err = exception.NoPermission
		}
	}
}
