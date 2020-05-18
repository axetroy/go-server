// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/rbac/accession"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

func GetAdminInfo(c helper.Context) (res schema.Response) {
	var (
		err  error
		data = schema.AdminProfile{}
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

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
		return
	}

	if err = mapstructure.Decode(adminInfo, &data.AdminProfilePure); err != nil {
		return
	}

	// 如果是超级管理员，则拥有全部权限
	if adminInfo.IsSuper == true {
		data.Accession = []string{}
		for _, v := range accession.AdminList {
			data.Accession = append(data.Accession, v.Name)
		}
	}

	if len(data.Accession) == 0 {
		data.Accession = []string{}
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetAdminInfoById(c helper.Context, adminId string) (res schema.Response) {
	var (
		err  error
		data = schema.AdminProfileWithToken{}
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

	tx = database.Db.Begin()

	handlerInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&handlerInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	adminInfo := model.Admin{
		Id: adminId,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
		return
	}

	if err = mapstructure.Decode(adminInfo, &data.AdminProfilePure); err != nil {
		return
	}

	if len(data.Accession) == 0 {
		data.Accession = []string{}
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var GetAdminInfoRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return GetAdminInfo(helper.NewContext(&c))
	})
})

var GetAdminInfoByIdRouter = router.Handler(func(c router.Context) {
	adminId := c.Param("admin_id")

	c.ResponseFunc(nil, func() schema.Response {
		return GetAdminInfoById(helper.NewContext(&c), adminId)
	})
})
