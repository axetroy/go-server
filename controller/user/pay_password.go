package user

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/password"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
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
		err error
		tx  *gorm.DB
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

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
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

	userInfo := model.User{Id: uid}

	tx = orm.DB.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if userInfo.PayPassword != nil {
		err = exception.PayPasswordSe
		return
	}

	// 更新交易密码
	if err = orm.DB.Model(userInfo).Update("pay_password", password.Generate(input.Password)).Error; err != nil {
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
		err error
		tx  *gorm.DB
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

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
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

	// TODO: 校验支付密码格式是否正确

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
		return
	}

	userInfo := model.User{Id: uid}

	tx = orm.DB.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 如果设置的密码和旧密码相同
	oldPwd := password.Generate(input.OldPassword)

	if userInfo.PayPassword != &oldPwd {
		err = exception.InvalidPassword
		return
	}

	newPwd := password.Generate(input.NewPassword)

	// 更新交易密码
	if err = orm.DB.Model(userInfo).Update("pay_password", newPwd).Error; err != nil {
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
