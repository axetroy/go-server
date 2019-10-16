// Copyright 2019 Axetroy. All rights reserved. MIT license.
package role

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func DeleteRoleByName(name string) {
	b := model.Role{}
	database.DeleteRowByTable(b.TableName(), "name", name)
}

func Delete(c controller.Context, roleName string) (res schema.Response) {
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

		helper.Response(&res, data, err)
	}()

	tx = database.Db.Begin()

	adminInfo := model.Admin{Id: c.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	roleInfo := model.Role{
		Name: roleName,
	}

	if err = tx.First(&roleInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AddressNotExist
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
		err = exception.RoleHadBeenUsed
		return
	}

	now := time.Now()
	timestamp := fmt.Sprintf("%v", now.UnixNano())

	// 我们重新更名这个角色，并且软删除
	if err = tx.Table(roleInfo.TableName()).Where("name = ? AND deleted_at IS NULL", roleInfo.Name).Update(model.Role{
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

func DeleteRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	roleName := c.Param("name")

	res = Delete(controller.NewContext(c), roleName)
}
