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
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/rbac/accession"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type UpdateParams struct {
	Description *string   `json:"description"`
	Accession   *[]string `json:"accession"`
	Note        *string   `json:"note"`
}

type UpdateUserRoleParams struct {
	Roles []string `json:"role"` // 要更新的用户角色
}

func Update(context schema.Context, roleName string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         role_schema.Role
		tx           *gorm.DB
		shouldUpdate bool
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
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
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

	roleInfo := role_model.Role{
		Name: roleName,
	}

	if err = tx.First(&roleInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrRoleNotExist
			return
		}
		return
	}

	updateModel := role_model.Role{}

	if input.Description != nil {
		shouldUpdate = true
		updateModel.Description = *input.Description
	}

	if input.Accession != nil {

		// 检验要更新的权限是否合法
		if accession.Valid(*input.Accession) == false {
			err = exception.ErrInvalidParams
			return
		}

		shouldUpdate = true
		updateModel.Accession = *input.Accession
	}

	if input.Note != nil {
		shouldUpdate = true
		updateModel.Note = input.Note
	}

	if shouldUpdate {
		if err = tx.Model(&roleInfo).UpdateColumns(&updateModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = ErrRoleNotExist
				return
			}
			return
		}
	}

	// 内建的角色是无法修改的
	if roleInfo.BuildIn == true {
		err = ErrRoleCannotUpdate
		return
	}

	if err = mapstructure.Decode(roleInfo, &data.RolePure); err != nil {
		return
	}

	data.CreatedAt = roleInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = roleInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateRouter(ctx *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	roleName := ctx.Param("name")

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = Update(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, roleName, input)
}

func UpdateUserRole(context schema.Context, userId string, input UpdateUserRoleParams) (res schema.Response) {
	var (
		err  error
		data user_schema.Profile
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
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

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

	userInfo := user_model.User{
		Id: userId,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	if len(input.Roles) > 20 {
		err = errors.New("一个用户不能拥有太多角色")
		return
	}

	// 确保要更新的角色存在
	for _, roleName := range input.Roles {
		roleInfo := role_model.Role{
			Name: roleName,
		}

		if err = tx.First(&roleInfo).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = ErrRoleNotExist
				return
			}
			return
		}
	}

	updateModel := user_model.User{
		Role: input.Roles,
	}

	if err = tx.Model(&userInfo).UpdateColumns(&updateModel).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateUserRoleRouter(ctx *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateUserRoleParams
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

	res = UpdateUserRole(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, userId, input)
}
