// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type UpdatePasswordParams struct {
	OldPassword string `json:"old_password" valid:"required~请输入旧密码"`
	NewPassword string `json:"new_password" valid:"required~请输入新密码"`
}

type UpdatePasswordByAdminParams struct {
	NewPassword string `json:"new_password" valid:"required~请输入新密码"`
}

func UpdatePassword(context schema.Context, input UpdatePasswordParams) (res schema.Response) {
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
			res.Data = false
		} else {
			res.Data = true
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
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

	tx = database.Db.Begin()

	userInfo := user_model.User{Id: context.Uid}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	// 验证密码是否正确
	if userInfo.Password != util.GeneratePassword(input.OldPassword) {
		err = exception.ErrInvalidPassword
		return
	}

	newPassword := util.GeneratePassword(input.NewPassword)

	if err = tx.Model(&userInfo).Update(user_model.User{Password: newPassword}).Error; err != nil {
		return
	}

	return
}

func UpdatePasswordByAdmin(context schema.Context, userId string, input UpdatePasswordByAdminParams) (res schema.Response) {
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
			res.Data = false
		} else {
			res.Data = true
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	// 检查是否是管理员
	adminInfo := admin_model.Admin{Id: context.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	// 只有超级管理员才能操作
	if adminInfo.IsSuper == false {
		err = exception.ErrNoPermission
		return
	}

	userInfo := user_model.User{Id: userId}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	newPassword := util.GeneratePassword(input.NewPassword)

	if err = tx.Model(&userInfo).Update(user_model.User{Password: newPassword}).Error; err != nil {
		return
	}

	return
}

func UpdatePasswordRouter(ctx *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdatePasswordParams
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

	res = UpdatePassword(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}

func UpdatePasswordByAdminRouter(ctx *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdatePasswordByAdminParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	userId := ctx.Param("user_id")

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = UpdatePasswordByAdmin(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, userId, input)
}
