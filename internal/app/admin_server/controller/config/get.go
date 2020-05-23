// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package config

import (
	"encoding/json"
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"time"
)

func Get(c helper.Context, configName string) (res schema.Response) {
	var (
		err  error
		data schema.Config
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

	// validator config name
	{
		c := model.Config{Name: configName}

		if err = c.IsValidConfigName(); err != nil {
			return
		}
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

	ConfigInfo := model.Config{
		Name: configName,
	}

	if err = tx.First(&ConfigInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = json.Unmarshal([]byte(ConfigInfo.Fields), &data.Fields); err != nil {
		err = exception.InvalidParams.New(err.Error())
		return
	}

	data.CreatedAt = ConfigInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = ConfigInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

var GetRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return Get(helper.NewContext(&c), c.Param("config_name"))
	})
})
