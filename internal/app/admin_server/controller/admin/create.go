// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/authentication"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type CreateAdminParams struct {
	Account  string `json:"account" validate:"required,min=1,max=36" comment:"帐号"`  // 管理员账号，登陆凭证
	Password string `json:"password" validate:"required,min=6,max=36" comment:"密码"` // 管理员密码
	Name     string `json:"name" validate:"required,min=2,max=36" comment:"名称"`     // 管理员名称，注册后不可修改
}

// 创建管理员
func CreateAdmin(input CreateAdminParams, isSuper bool) (res schema.Response) {
	var (
		err  error
		data schema.AdminProfileWithToken
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

		helper.Response(&res, data, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	n := model.Admin{Username: input.Account}

	if err = tx.Where(&n).First(&n).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = nil
		} else {
			err = exception.Database.New(err.Error())
			return
		}
	} else {
		err = exception.AdminExist
		return
	}

	adminInfo := model.Admin{
		Username:  input.Account,
		Name:      input.Name,
		Password:  util.GeneratePassword(input.Password),
		Status:    model.AdminStatusInit,
		Accession: []string{},
		IsSuper:   isSuper,
	}

	if err = tx.Create(&adminInfo).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(adminInfo, &data.AdminProfilePure); err != nil {
		return
	}

	// generate token
	if data.Token, err = authentication.Gateway(true).Generate(adminInfo.Id); err != nil {
		return
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var CreateAdminRouter = router.Handler(func(c router.Context) {
	var (
		input CreateAdminParams
		err   error
	)

	defer func() {
		c.ResponseFunc(err, func() schema.Response {
			return CreateAdmin(input, false)
		})
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		return
	}

	adminInfo := model.Admin{
		Id: c.Uid(),
	}

	if err = database.Db.Where(&adminInfo).First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
		return
	}

	if !adminInfo.IsSuper {
		err = exception.AdminNotSuper
		return
	}

})
