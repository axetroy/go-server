package middleware

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"gitlab.com/axetroy/server/exception"
	"gitlab.com/axetroy/server/model"
	"gitlab.com/axetroy/server/orm"
	"gitlab.com/axetroy/server/response"
	"gitlab.com/axetroy/server/services/password"
	"net/http"
)

// 验证交易密码的中间价
func AuthPayPassword(context *gin.Context) {
	var (
		err     error
		session *xorm.Session
		tx      bool
	)

	defer func() {

		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.Unknown
			}
		}

		if tx {
			if err != nil {
				_ = session.Rollback()
			} else {
				err = session.Commit()
			}
		}

		if session != nil {
			session.Close()
		}

		// 如果有报错的话，那么不会进入到路由中
		if err != nil {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})

			context.Abort()

			return
		}
	}()

	payPassword := context.GetHeader("X-Pay-Password")

	if len(payPassword) == 0 {
		err = exception.RequirePayPassword
		return
	}

	uid := context.GetInt64("uid")

	if uid == 0 {
		err = exception.UserNotLogin
		return
	}

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	user := model.User{Id: uid}

	var isExist bool

	if isExist, err = session.Get(&user); err != nil {
		return
	}

	if isExist == false {
		err = exception.UserNotExist
		return
	}

	// 校验密码是否正确
	if *user.PayPassword != password.Generate(payPassword) {
		err = exception.InvalidPassword
		return
	}

}
