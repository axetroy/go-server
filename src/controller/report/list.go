// Copyright 2019 Axetroy. All rights reserved. MIT license.
package report

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	Type   *model.ReportType   `json:"type" form:"type"`     // 类型
	Status *model.ReportStatus `json:"status" form:"status"` // 状态
}

type QueryAdmin struct {
	Query
	Uid string `json:"uid"`
}

func GetList(c controller.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]schema.Report, 0)
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

		helper.ResponseList(&res, data, meta, err)
	}()

	query := input.Query

	query.Normalize()

	list := make([]model.Report, 0)

	filter := map[string]interface{}{}

	filter["uid"] = c.Uid

	if input.Type != nil {
		filter["type"] = *input.Type
	}

	if input.Status != nil {
		filter["status"] = *input.Status
	}

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Where(filter).Find(&list).Error; err != nil {
		return
	}

	var total int64

	if err = database.Db.Model(&model.Report{}).Where(filter).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.Report{}
		if er := mapstructure.Decode(v, &d.ReportPure); er != nil {
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

func GetListRouter(c *gin.Context) {
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

	res = GetList(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, input)
}

func GetListByAdmin(c controller.Context, input QueryAdmin) (res schema.List) {
	var (
		err  error
		data = make([]schema.Report, 0)
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

		helper.ResponseList(&res, data, meta, err)
	}()

	query := input.Query

	query.Normalize()

	list := make([]model.Report, 0)

	filter := map[string]interface{}{}

	if input.Type != nil {
		filter["type"] = *input.Type
	}

	if input.Status != nil {
		filter["status"] = *input.Status
	}

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Where(filter).Find(&list).Error; err != nil {
		return
	}

	var total int64

	for _, v := range list {
		d := schema.Report{}
		if er := mapstructure.Decode(v, &d.ReportPure); er != nil {
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

func GetListByAdminRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.List{}
		input QueryAdmin
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

	res = GetListByAdmin(controller.NewContext(c), input)
}
