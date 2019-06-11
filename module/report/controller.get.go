// Copyright 2019 Axetroy. All rights reserved. MIT license.
package report

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/report/report_model"
	"github.com/axetroy/go-server/module/report/report_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func GetReportByUser(context schema.Context, id string) (res schema.Response) {
	var (
		err  error
		data = report_schema.Report{}
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	reportInfo := report_model.Report{
		Id:  id,
		Uid: context.Uid,
	}

	if err = database.Db.First(&reportInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrNoData
		}
		return
	}

	if err = mapstructure.Decode(reportInfo, &data.ReportPure); err != nil {
		return
	}

	data.CreatedAt = reportInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = reportInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetReportRouter(ctx *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	id := ctx.Param("report_id")

	res = GetReportByUser(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}

func GetReportByAdmin(context schema.Context, id string) (res schema.Response) {
	var (
		err  error
		data = report_schema.Report{}
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	reportInfo := report_model.Report{
		Id: id,
	}

	if err = database.Db.First(&reportInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrNoData
		}
		return
	}

	if err = mapstructure.Decode(reportInfo, &data.ReportPure); err != nil {
		return
	}

	data.CreatedAt = reportInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = reportInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetReportByAdminRouter(ctx *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	id := ctx.Param("report_id")

	res = GetReportByAdmin(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}
