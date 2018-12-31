package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/password"
	"net/http"
	"strconv"
)

type UpdatePasswordParams struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func UpdatePassword(context *gin.Context) {
	var (
		err     error
		uid     int64
		input   UpdatePasswordParams
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
				Message: "更新成功",
				Data:    true,
			})
		}
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		return
	}

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
		return
	}

	// TODO 校验input是否正确

	if val, isExist := context.Get("uid"); isExist != true {

	} else {
		if uid, err = strconv.ParseInt(fmt.Sprintf("%v", val), 10, 64); err != nil {
			return
		}
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
}
