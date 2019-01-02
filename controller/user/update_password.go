package user

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/password"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"net/http"
)

type UpdatePasswordParams struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UpdatePassword(uid string, input UpdatePasswordParams) (res response.Response) {
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
			res.Data = nil
			res.Data = false
		} else {
			res.Data = true
			res.Message = "更新成功"
			res.Status = response.StatusSuccess
		}
	}()

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
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

	if user.Password != password.Generate(input.OldPassword) {
		err = exception.InvalidPassword
		return
	}

	user.Password = password.Generate(input.NewPassword)
	if _, err = session.Cols("password").Update(&user); err == nil {
		return
	}
	return
}

func UpdatePasswordRouter(context *gin.Context) {
	var (
		err   error
		res   = response.Response{}
		input UpdatePasswordParams
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

	res = UpdatePassword(context.GetString("uid"), input)
}
