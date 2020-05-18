// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package role

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

// 获取角色信息
func Get(roleName string) (res schema.Response) {
	var (
		err  error
		data = schema.Role{}
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

		helper.Response(&res, data, nil, err)
	}()

	roleInfo := model.Role{
		Name: roleName,
	}

	if err = database.Db.First(&roleInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.RoleNotExist
		}
		return
	}

	if err = mapstructure.Decode(roleInfo, &data.RolePure); err != nil {
		return
	}

	data.CreatedAt = roleInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = roleInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var GetRouter = router.Handler(func(c router.Context) {
	roleName := c.Param("name")

	c.ResponseFunc(nil, func() schema.Response {
		return Get(roleName)
	})
})
