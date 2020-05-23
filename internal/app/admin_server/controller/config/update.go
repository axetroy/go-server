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

func Update(c helper.Context, configName string, fields []byte) (res schema.Response) {
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
		c := model.Config{
			Name:   configName,
			Fields: string(fields),
		}

		if err = c.IsValidConfigName(); err != nil {
			return
		}

		if err = c.IsValidConfigField(); err != nil {
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

	configInfo := model.Config{
		Name: configName,
	}

	if err = tx.Model(&configInfo).Where(&configInfo).First(&configInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	configInfo.Fields = string(fields)

	if err = tx.Model(&configInfo).Where(model.Config{Name: configName}).Update(&configInfo).Error; err != nil {
		return
	}

	if err = json.Unmarshal(fields, &data.Fields); err != nil {
		err = exception.InvalidParams.New(err.Error())
		return
	}

	data.CreatedAt = configInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = configInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

var UpdateRouter = router.Handler(func(c router.Context) {
	var body []byte

	body, err := c.GetBody()

	c.ResponseFunc(err, func() schema.Response {
		return Update(helper.NewContext(&c), c.Param("config_name"), body)
	})
})
