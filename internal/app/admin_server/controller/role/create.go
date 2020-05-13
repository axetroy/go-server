// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package role

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/rbac/accession"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
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

func Create(c helper.Context, input CreateParams) (res schema.Response) {
	var (
		err  error
		data schema.Role
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

	adminInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	if !adminInfo.IsSuper {
		err = exception.AdminNotSuper
		return
	}

	if accession.Valid(input.Accession) == false {
		err = exception.InvalidParams
		return
	}

	roleInfo := model.Role{
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

func CreateRouter(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Create(helper.NewContext(c), input)
}
