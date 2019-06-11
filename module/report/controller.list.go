// Copyright 2019 Axetroy. All rights reserved. MIT license.
package report

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/report/report_model"
	"github.com/axetroy/go-server/module/report/report_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	Type   *report_model.ReportType   `json:"type" form:"type"`     // 类型
	Status *report_model.ReportStatus `json:"status" form:"status"` // 状态
}

type QueryAdmin struct {
	Query
	Uid string `json:"uid"`
}

func GetList(context schema.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]report_schema.Report, 0)
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
				err = exception.ErrUnknown
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

	list := make([]report_model.Report, 0)

	var total int64

	search := report_model.Report{
		Uid: context.Uid,
	}

	if input.Type != nil {
		search.Type = *input.Type
	}

	if input.Status != nil {
		search.Status = *input.Status
	}

	if err = database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Where(&search).Find(&list).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := report_schema.Report{}
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
		err = exception.ErrInvalidParams
		return
	}

	res = GetList(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}

func GetListByAdmin(context schema.Context, input QueryAdmin) (res schema.List) {
	var (
		err  error
		data = make([]report_schema.Report, 0)
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
				err = exception.ErrUnknown
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

	list := make([]report_model.Report, 0)

	var total int64

	search := report_model.Report{}

	if input.Type != nil {
		search.Type = *input.Type
	}

	if input.Status != nil {
		search.Status = *input.Status
	}

	if err = database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Where(&search).Find(&list).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := report_schema.Report{}
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

	return
}

func GetListByAdminRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindQuery(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = GetListByAdmin(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
