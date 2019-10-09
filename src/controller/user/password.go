// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/src/validator"
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

func UpdatePassword(context controller.Context, input UpdatePasswordParams) (res schema.Response) {
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

		helper.Response(&res, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
		return
	}

	tx = database.Db.Begin()

	userInfo := model.User{Id: context.Uid}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 验证密码是否正确
	if userInfo.Password != util.GeneratePassword(input.OldPassword) {
		err = exception.InvalidPassword
		return
	}

	newPassword := util.GeneratePassword(input.NewPassword)

	if err = tx.Model(&userInfo).Update(model.User{Password: newPassword}).Error; err != nil {
		return
	}

	return
}

func UpdatePasswordByAdmin(context controller.Context, userId string, input UpdatePasswordByAdminParams) (res schema.Response) {
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

		helper.Response(&res, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	// 检查是否是管理员
	adminInfo := model.Admin{Id: context.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	// 只有超级管理员才能操作
	if adminInfo.IsSuper == false {
		err = exception.NoPermission
		return
	}

	userInfo := model.User{Id: userId}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	newPassword := util.GeneratePassword(input.NewPassword)

	if err = tx.Model(&userInfo).Update(model.User{Password: newPassword}).Error; err != nil {
		return
	}

	return
}

func UpdatePasswordRouter(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = UpdatePassword(controller.NewContext(c), input)
}

func UpdatePasswordByAdminRouter(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	}()

	userId := c.Param("user_id")

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = UpdatePasswordByAdmin(controller.NewContext(c), userId, input)
}
