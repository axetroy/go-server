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

func SetPayPassword(uid string, input SetPayPasswordParams) (res response.Response) {
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
			res.Status = response.StatusSuccess
		}
	}()

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

	return
}

func SetPayPasswordRouter(context *gin.Context) {
	var (
		err   error
		res   = response.Response{}
		input SetPayPasswordParams
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

	res = SetPayPassword(context.GetString("uid"), input)
}

func UpdatePayPassword(uid string, input UpdatePayPasswordParams) (res response.Response) {
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
			res.Status = response.StatusSuccess
		}
	}()

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
		return
	}

	// TODO: 校验支付密码格式是否正确

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

	return
}

func UpdatePayPasswordRouter(context *gin.Context) {
	var (
		err   error
		res   = response.Response{}
		input UpdatePayPasswordParams
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

	res = UpdatePayPassword(context.GetString("uid"), input)
}
