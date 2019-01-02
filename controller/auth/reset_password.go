package auth

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/password"
	"github.com/axetroy/go-server/services/redis"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"net/http"
)

type ResetPasswordParams struct {
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func ResetPassword(input ResetPasswordParams) (res response.Response) {
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

		if err != nil {
			res.Message = err.Error()
			res.Data = false
		} else {
			res.Status = response.StatusSuccess
			res.Data = true
		}
	}()

	var (
		uid string
	)

	if uid, err = redis.ResetCode.Get(input.Code).Result(); err != nil {
		err = exception.InvalidResetCode
		return
	}

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	defer func() {
		if err != nil {
			_ = session.Rollback()
		} else {
			_ = session.Commit()
		}
	}()

	user := model.User{Id: uid}

	var isExist bool

	if isExist, err = session.Get(&user); err != nil {
		return
	}

	if isExist == false {
		err = exception.UserNotExist
		return
	}

	user.Password = password.Generate(input.NewPassword)

	if _, err = session.Cols("password").Update(&user); err != nil {
		return
	}

	// delete reset code from redis
	if err = redis.ResetCode.Del(input.Code).Err(); err != nil {
		return
	}

	// TODO: 安全起见，发送一封邮件/短信告知用户
	return
}

func ResetPasswordRouter(context *gin.Context) {
	var (
		input ResetPasswordParams
		err   error
		res   = response.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = ResetPassword(input)
}
