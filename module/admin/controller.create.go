// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/admin/admin_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateAdminParams struct {
	Account  string `json:"account" valid:"required~请输入管理员账号"`  // 管理员账号，登陆凭证
	Password string `json:"password" valid:"required~请输入管理员密码"` // 管理员密码
	Name     string `json:"name" valid:"required~请输入管理员名称"`     // 管理员名称，注册后不可修改
}

// 创建管理员
func CreateAdmin(input CreateAdminParams, isSuper bool) (res schema.Response) {
	var (
		err          error
		data         admin_schema.AdminProfile
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
				err = common_error.ErrUnknown
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}

	}()

	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = common_error.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	n := admin_model.Admin{Username: input.Account}

	if tx.Where(&n).First(&n).RecordNotFound() == false {
		err = ErrAdminExist
		return
	}

	adminInfo := admin_model.Admin{
		Username:  input.Account,
		Name:      input.Name,
		Password:  util.GeneratePassword(input.Password),
		Status:    admin_model.AdminStatusInit,
		Accession: []string{},
		IsSuper:   isSuper,
	}

	if err = tx.Create(&adminInfo).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(adminInfo, &data.AdminProfilePure); err != nil {
		return
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func CreateAdminRouter(ctx *gin.Context) {
	var (
		input CreateAdminParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	uid := ctx.GetString(middleware.ContextUidField)

	adminInfo := admin_model.Admin{
		Id: uid,
	}

	if err = database.Db.Where(&adminInfo).First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrAdminNotExist
			return
		}
		return
	}

	if adminInfo.IsSuper == false {
		err = ErrAdminNotSuper
		return
	}

	res = CreateAdmin(input, false)
}
