// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/address"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/role/role_model"
	"github.com/axetroy/go-server/module/role/role_schema"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func DeleteRoleByName(name string) {
	b := role_model.Role{}
	database.DeleteRowByTable(b.TableName(), "name", name)
}

func Delete(context schema.Context, roleName string) (res schema.Response) {
	var (
		err  error
		data role_schema.Role
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

	adminInfo := admin_model.Admin{Id: context.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	roleInfo := role_model.Role{
		Name: roleName,
	}

	if err = tx.First(&roleInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = address.ErrAddressNotExist
			return
		}
		return
	}

	// 查询是否有用户属于这个角色，如果有，不允许删除
	var roleUsersNum int64

	if err = tx.Raw(fmt.Sprintf(`SELECT COUNT(*) FROM "user"  WHERE "user"."role" IN ('{%s}')`, roleInfo.Name)).Count(&roleUsersNum).Error; err != nil {
		return
	}

	if roleUsersNum > 0 {
		err = ErrRoleHadBeenUsed
		return
	}

	now := time.Now()
	timestamp := fmt.Sprintf("%v", now.UnixNano())

	// 我们重新更名这个角色，并且软删除
	if err = tx.Table(roleInfo.TableName()).Where("name = ? AND deleted_at IS NULL", roleInfo.Name).Update(role_model.Role{
		Name:      roleInfo.Name + "_" + timestamp,
		DeletedAt: &now,
	}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(roleInfo, &data.RolePure); err != nil {
		return
	}

	data.CreatedAt = roleInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = roleInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func DeleteRouter(ctx *gin.Context) {
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

	roleName := ctx.Param("name")

	res = Delete(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, roleName)
}
