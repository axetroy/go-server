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

type SetPayPasswordParams struct {
	Password string `json:"password"`
}

type UpdatePayPasswordParams struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func SetPayPassword(context *gin.Context) {
	var (
		err     error
		uid     string
		input   SetPayPasswordParams
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
				Message: "设置成功",
				Data:    true,
			})
		}
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		return
	}

	// TODO 校验input是否正确

	uid = context.GetString("uid")

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

	if user.PayPassword != nil {
		err = errors.New("交易密码已设置")
		return
	}

	pwd := password.Generate(input.Password)

	user.PayPassword = &pwd

	if _, err = session.Cols("pay_password").Update(&user); err == nil {
		return
	}
}

func UpdatePayPassword(context *gin.Context) {
	var (
		err     error
		uid     string
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

	// TODO: 校验支付密码格式是否正确

	uid = context.GetString("uid")

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

	// 如果设置的密码和旧密码相同
	oldPwd := password.Generate(input.OldPassword)

	if user.PayPassword != &oldPwd {
		err = exception.InvalidPassword
		return
	}

	newPwd := password.Generate(input.NewPassword)

	user.PayPassword = &newPwd
	if _, err = session.Cols("pay_password").Update(&user); err == nil {
		return
	}
}
