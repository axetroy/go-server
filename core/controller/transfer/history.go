// Copyright 2019 Axetroy. All rights reserved. MIT license.
package transfer

import (
	"errors"
	"github.com/axetroy/go-server/core/controller"
	"github.com/axetroy/go-server/core/exception"
	"github.com/axetroy/go-server/core/helper"
	"github.com/axetroy/go-server/core/middleware"
	"github.com/axetroy/go-server/core/model"
	"github.com/axetroy/go-server/core/schema"
	"github.com/axetroy/go-server/core/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
}

func GetHistory(c controller.Context, input Query) (res schema.List) {
	var (
		err  error
		tx   *gorm.DB
		data = make([]schema.TransferLog, 0)
		meta = &schema.Meta{}
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

		helper.ResponseList(&res, data, meta, err)
	}()

	tx = database.Db.Begin()

	userInfo := model.User{Id: c.Uid}

	if err = tx.Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	query := input.Query

	query.Normalize()

	list := make([]model.TransferLog, 0)

	condition := QueryParams{
		From: &c.Uid,
	}

	// 联表查询
	countSQL := GenerateTransferLogSQL(condition, query.Limit, true)
	listSQL := GenerateTransferLogSQL(condition, query.Limit, false)

	var total int64

	if err = tx.Raw(countSQL).Count(&total).Error; err != nil {
		return
	}

	if err = tx.Raw(listSQL).Scan(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.TransferLog{}
		if er := mapstructure.Decode(v, &d.TransferLogPure); er != nil {
			err = er
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(list)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}

func GetHistoryRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.List{}
		input Query
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = GetHistory(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, input)
}
