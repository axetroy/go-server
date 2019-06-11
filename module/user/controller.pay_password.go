// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/service/email"
	"github.com/axetroy/go-server/service/redis"
	"github.com/axetroy/go-server/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type SetPayPasswordParams struct {
	Password        string `json:"password" valid:"required~请输入密码,int~请输入纯数字的密码,length(6|6)~密码长度为6位"`
	PasswordConfirm string `json:"password_confirm" valid:"required~请输入确认密码,int~请输入纯数字的确认密码,length(6|6)~确认密码长度为6位"`
}

type UpdatePayPasswordParams struct {
	OldPassword string `json:"old_password" valid:"required~请输入旧密码,int~请输入纯数字的旧密码,length(6|6)~旧密码长度为6位"`
	NewPassword string `json:"new_password" valid:"required~请输入新密码,int~请输入纯数字的新密码,length(6|6)~新密码长度为6位"`
}

type ResetPayPasswordParams struct {
	Code        string `json:"code" valid:"required~请输入重置码"`                                                // 重置码
	NewPassword string `json:"new_password" valid:"required~请输入新的交易密码,int~请输入纯数字的旧密码,length(6|6)~新密码长度为6位"` // 新的交易密码
}

func GenerateResetPayPasswordCode(uid string) string {
	codeId := "reset-pay-" + util.GenerateId() + uid
	return util.MD5(codeId)
}

func SetPayPassword(context schema.Context, input SetPayPasswordParams) (res schema.Response) {
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
				err = exception.ErrUnknown
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
		err = exception.ErrInvalidParams
		return
	}

	if input.Password != input.PasswordConfirm {
		err = user_error.ErrInvalidConfirmPassword
		return
	}

	userInfo := user_model.User{Id: context.Uid}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	if userInfo.PayPassword != nil {
		err = exception.ErrPayPasswordSet
		return
	}

	newPassword := util.GeneratePassword(input.Password)

	// 更新交易密码
	if err = database.Db.Model(userInfo).Update("pay_password", newPassword).Error; err != nil {
		return
	}

	return
}

func SetPayPasswordRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = SetPayPassword(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}

func UpdatePayPassword(context schema.Context, input UpdatePayPasswordParams) (res schema.Response) {
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
				err = exception.ErrUnknown
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
		err = exception.ErrInvalidParams
		return
	}

	if input.OldPassword == input.NewPassword {
		err = exception.ErrPasswordDuplicate
		return
	}

	userInfo := user_model.User{Id: context.Uid}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	// 如果还没有设置过交易密码，就不会有更新
	if userInfo.PayPassword == nil {
		err = user_error.ErrErrRequireErrPayPasswordSet
		return
	}

	oldPwd := util.GeneratePassword(input.OldPassword)

	// 旧密码不匹配
	if *userInfo.PayPassword != oldPwd {
		err = exception.ErrInvalidPassword
		return
	}

	newPwd := util.GeneratePassword(input.NewPassword)

	// 更新交易密码
	if err = database.Db.Model(userInfo).Update("pay_password", newPwd).Error; err != nil {
		return
	}

	return
}

func UpdatePayPasswordRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = UpdatePayPassword(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}

func SendResetPayPassword(context schema.Context) (res schema.Response) {
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
				err = exception.ErrUnknown
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

	userInfo := user_model.User{Id: context.Uid}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	// 生成重置码
	var resetCode = GenerateResetPayPasswordCode(userInfo.Id)

	// redis缓存重置码
	if err = redis.ResetCodeClient.Set(resetCode, userInfo.Id, time.Minute*10).Err(); err != nil {
		return
	}

	if userInfo.Email != nil {
		// 发送邮件
		go func() {
			e := email.NewMailer()
			_ = e.SendForgotTradePasswordEmail(*userInfo.Email, resetCode)
		}()
	} else if userInfo.Phone != nil {
		// TODO: 发送手机验证码
		go func() {

		}()
	} else {
		// 无效的用户
		err = exception.ErrNoData
	}

	return
}

func SendResetPayPasswordRouter(ctx *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	res = SendResetPayPassword(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	})
}

func ResetPayPassword(context schema.Context, input ResetPayPasswordParams) (res schema.Response) {
	var (
		err          error
		tx           *gorm.DB
		isValidInput bool
		uid          string
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.ErrUnknown
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
		err = exception.ErrInvalidParams
		return
	}

	userInfo := user_model.User{Id: context.Uid}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	// 如果还没有设置过交易密码，就不会有重置
	if userInfo.PayPassword == nil {
		err = user_error.ErrErrRequireErrPayPasswordSet
		return
	}

	if uid, err = redis.ResetCodeClient.Get(input.Code).Result(); err != nil {
		err = user_error.ErrInvalidResetCode
		return
	}

	// 即使有了重置码，不是自己的账号也不能用
	if userInfo.Id != uid {
		err = exception.ErrNoPermission
		return
	}

	// 更新交易密码
	if err = database.Db.Model(userInfo).Update("pay_password", input.NewPassword).Error; err != nil {
		return
	}

	// 重置密码之后，删除重置码
	if _, err = redis.ResetCodeClient.Del(input.Code).Result(); err != nil {
		return
	}

	return
}

func ResetPayPasswordRouter(ctx *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input ResetPayPasswordParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = ResetPayPassword(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
