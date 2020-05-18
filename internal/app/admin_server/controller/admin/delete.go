// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin

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

func DeleteAdminByAccount(account string) {
	database.DeleteRowByTable("admin", "username", account)
}

func DeleteAdminById(c helper.Context, adminId string) (res schema.Response) {
	var (
		err  error
		data schema.AdminProfile
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

	targetAdminInfo := model.Admin{
		Id: adminId,
	}

	if err = tx.First(&targetAdminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
		return
	}

	myInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&myInfo).Error; err != nil {
		return
	}

	// 超级管理员才能操作
	if myInfo.IsSuper == false {
		err = exception.NoPermission
		return
	}

	if err = tx.Delete(model.Admin{
		Id:      targetAdminInfo.Id,
		IsSuper: false, // 超级管理员无法被删除
	}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(targetAdminInfo, &data.AdminProfilePure); err != nil {
		return
	}

	if len(data.Accession) == 0 {
		data.Accession = []string{}
	}

	data.CreatedAt = targetAdminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = targetAdminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var DeleteAdminByIdRouter = router.Handler(func(c router.Context) {
	id := c.Param("admin_id")

	c.ResponseFunc(nil, func() schema.Response {
		return DeleteAdminById(helper.NewContext(&c), id)
	})
})
