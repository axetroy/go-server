package auth

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"gitlab.com/axetroy/server/exception"
	"gitlab.com/axetroy/server/model"
	"gitlab.com/axetroy/server/orm"
	"gitlab.com/axetroy/server/response"
	"gitlab.com/axetroy/server/services/password"
	"gitlab.com/axetroy/server/services/redis"
	"net/http"
	"strconv"
)

type ResetPasswordParams struct {
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func ResetPassword(context *gin.Context) {
	var (
		input   ResetPasswordParams
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
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})
		} else {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusSuccess,
				Message: "密码重置成功",
				Data:    true,
			})
		}
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	// TODO: 参数校验

	var (
		uidStr string
		uid    int64
	)

	if uidStr, err = redis.ResetCode.Get(input.Code).Result(); err != nil {
		err = exception.InvalidResetCode
		return
	}

	if uid, err = strconv.ParseInt(uidStr, 10, 64); err != nil {
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

}
