// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/role/role_model"
	"github.com/axetroy/go-server/module/role/role_schema"
	"github.com/axetroy/go-server/rbac/accession"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateParams struct {
	Name        string   `json:"name" valid:"required~请输入角色名"`       // 角色名
	Description string   `json:"description" valid:"required~请输入描述"` // 描述
	Accession   []string `json:"accession" valid:"required~请输入权限"`   // 权限列表
	Note        *string  `json:"note"`                               // 备注
}

func Create(context schema.Context, input CreateParams) (res schema.Response) {
	var (
		err          error
		data         role_schema.Role
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
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

	adminInfo := admin_model.Admin{
		Id: context.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	if !adminInfo.IsSuper {
		err = admin.ErrAdminNotSuper
		return
	}

	if accession.Valid(input.Accession) == false {
		err = exception.ErrInvalidParams
		return
	}

	roleInfo := role_model.Role{
		Name:        input.Name,
		Description: input.Description,
		Accession:   input.Accession,
	}

	if err = tx.Create(&roleInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(roleInfo, &data.RolePure); er != nil {
		err = er
		return
	}

	data.CreatedAt = roleInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = roleInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func CreateRouter(ctx *gin.Context) {
	var (
		input CreateParams
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
		err = exception.ErrInvalidParams
		return
	}

	res = Create(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
