// Copyright 2019 Axetroy. All rights reserved. MIT license.
package login

import (
	"errors"
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

func GetLatestLoginLog(c controller.Context) (res schema.Response) {
	var (
		err  error
		data = schema.LogLogin{}
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

		helper.Response(&res, data, err)
	}()

	logInfo := model.LoginLog{
		Uid: c.Uid,
	}

	query := schema.Query{}

	query.Normalize()

	if err = query.Order(database.Db.Where(&logInfo).Preload("User")).First(&logInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(logInfo, &data.LogLoginPure); err != nil {
		return
	}

	if err = mapstructure.Decode(logInfo.User, &data.User); err != nil {
		return
	}

	data.CreatedAt = logInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = logInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetLoginLog(id string) (res schema.Response) {
	var (
		err  error
		data = schema.LogLogin{}
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

		helper.Response(&res, data, err)
	}()

	logInfo := model.LoginLog{
		Id: id,
	}

	query := schema.Query{}

	query.Normalize()

	if err = query.Order(database.Db.Where(&logInfo).Preload("User")).First(&logInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(logInfo, &data.LogLoginPure); err != nil {
		return
	}

	if err = mapstructure.Decode(logInfo.User, &data.User); err != nil {
		return
	}

	data.CreatedAt = logInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = logInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetLoginLogRouter(c *gin.Context) {
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

	id := c.Param("log_id")

	res = GetLoginLog(id)
}
