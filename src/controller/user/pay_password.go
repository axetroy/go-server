package user

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/axetroy/go-server/src/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type SetPayPasswordParams struct {
	Password        string `json:"password" valid:"required~请输入密码,int~请输入纯数字的密码,length(6|6)~密码长度为6位"`
	PasswordConfirm string `json:"password_confirm" valid:"required~请输入确认密码,int~请输入纯数字的确认密码,length(6|6)~确认密码长度为6位"`
}

type UpdatePayPasswordParams struct {
	OldPassword string `json:"old_password" valid:"required~请输入旧密码,int~请输入纯数字的旧密码,length(6|6)~旧密码长度为6位"`
	NewPassword string `json:"new_password" valid:"required~请输入新密码,int~请输入纯数字的新密码,length(6|6)~新密码长度为6位"`
}

func SetPayPassword(context controller.Context, input SetPayPasswordParams) (res schema.Response) {
	var (
		err          error
		tx           *gorm.DB
		isValidInput bool
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
			res.Status = schema.StatusSuccess
		}
	}()

	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	if input.Password != input.PasswordConfirm {
		err = exception.InvalidConfirmPassword
		return
	}

	userInfo := model.User{Id: context.Uid}

	tx = service.Db.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if userInfo.PayPassword != nil {
		err = exception.PayPasswordSet
		return
	}

	newPassword := util.GeneratePassword(input.Password)

	// 更新交易密码
	if err = service.Db.Model(userInfo).Update("pay_password", newPassword).Error; err != nil {
		return
	}

	return
}

func SetPayPasswordRouter(context *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
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

	res = SetPayPassword(controller.Context{
		Uid: context.GetString("uid"),
	}, input)
}

func UpdatePayPassword(context controller.Context, input UpdatePayPasswordParams) (res schema.Response) {
	var (
		err          error
		tx           *gorm.DB
		isValidInput bool
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
			res.Status = schema.StatusSuccess
		}
	}()

	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
		return
	}

	userInfo := model.User{Id: context.Uid}

	tx = service.Db.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 如果设置的密码和旧密码相同
	oldPwd := util.GeneratePassword(input.OldPassword)

	if userInfo.PayPassword != &oldPwd {
		err = exception.InvalidPassword
		return
	}

	newPwd := util.GeneratePassword(input.NewPassword)

	// 更新交易密码
	if err = service.Db.Model(userInfo).Update("pay_password", newPwd).Error; err != nil {
		return
	}

	return
}

func UpdatePayPasswordRouter(context *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
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

	res = UpdatePayPassword(controller.Context{
		Uid: context.GetString("uid"),
	}, input)
}
