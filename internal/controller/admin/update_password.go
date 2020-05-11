// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin

import (
	"errors"
	"github.com/axetroy/go-server/internal/helper"
	"github.com/axetroy/go-server/internal/validator"
	"net/http"
	"time"

	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type UpdatePasswordParams struct {
	OldPassword     string `json:"old_password" valid:"required~请输入旧密码"`      // 旧密码
	NewPassword     string `json:"new_password" valid:"required~请输入新密码"`      // 新密码
	ConfirmPassword string `json:"confirm_password" valid:"required~请输入确认密码"` // 确认密码
}

func UpdatePassword(c controller.Context, input UpdatePasswordParams) (res schema.Response) {
	var (
		err  error
		data schema.AdminProfile
		tx   *gorm.DB
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if input.NewPassword != input.ConfirmPassword {
		err = exception.InvalidConfirmPassword
		return
	}

	tx = database.Db.Begin()

	myInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&myInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	// 校验密码是否正确
	if myInfo.Password != util.GeneratePassword(input.OldPassword) {
		err = exception.InvalidOldPassword
		return
	}

	newPassword := util.GeneratePassword(input.NewPassword)

	if err = tx.Model(&myInfo).Update("password", newPassword).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
		return
	}

	if err = mapstructure.Decode(myInfo, &data.AdminProfilePure); err != nil {
		return
	}

	if len(data.Accession) == 0 {
		data.Accession = []string{}
	}

	data.CreatedAt = myInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = myInfo.UpdatedAt.Format(time.RFC3339Nano)

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

	res = UpdatePassword(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, input)
}
