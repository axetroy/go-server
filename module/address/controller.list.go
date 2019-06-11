// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/address/address_model"
	"github.com/axetroy/go-server/module/address/address_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	//Status model.NewsStatus `json:"status" form:"status"`
}

func GetList(context schema.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]address_schema.Address, 0) // 输出到外部的结果
		list = make([]address_model.Address, 0)  // 数据库查询出来的原始结果
		meta = &schema.Meta{}
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
				err = common_error.ErrUnknown
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
			res.Meta = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
			res.Meta = meta
		}
	}()

	query := input.Query

	query.Normalize()

	var total int64

	tx = database.Db.Begin()

	if err = tx.Limit(query.Limit).Offset(query.Limit*query.Page).Where(address_model.Address{Uid: context.Uid}, context.Uid).Find(&list).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := address_schema.Address{}
		if er := mapstructure.Decode(v, &d.AddressPure); er != nil {
			err = er
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit

	return
}

func GetListRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindQuery(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	res = GetList(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
