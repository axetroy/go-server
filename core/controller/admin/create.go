// Copyright 2019 Axetroy. All rights reserved. MIT license.
package admin

import (
	"errors"
	"github.com/axetroy/go-server/core/exception"
	"github.com/axetroy/go-server/core/helper"
	"github.com/axetroy/go-server/core/middleware"
	"github.com/axetroy/go-server/core/model"
	"github.com/axetroy/go-server/core/schema"
	"github.com/axetroy/go-server/core/service/database"
	"github.com/axetroy/go-server/core/service/token"
	"github.com/axetroy/go-server/core/util"
	"github.com/axetroy/go-server/core/validator"
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	n := model.Admin{Username: input.Account}

	if tx.Where(&n).First(&n).RecordNotFound() == false {
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
	if data.Token, err = token.Generate(adminInfo.Id, true); err != nil {
		return
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func CreateAdminRouter(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	uid := c.GetString(middleware.ContextUidField)

	adminInfo := model.Admin{
		Id: uid,
	}

	if err = database.Db.Where(&adminInfo).First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
		return
	}

	if adminInfo.IsSuper == false {
		err = exception.AdminNotSuper
		return
	}

	res = CreateAdmin(input, false)
}
