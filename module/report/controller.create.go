// Copyright 2019 Axetroy. All rights reserved. MIT license.
package report

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/exception"
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

type CreateParams struct {
	Title       string                  `json:"title" valid:"required~请填写标题"`     // 标题
	Content     string                  `json:"content" valid:"required~请填写反馈内容"` // 内容
	Type        report_model.ReportType `json:"type" valid:"required~请填写反馈类型"`    // 反馈类型
	Screenshots []string                `json:"screenshots"`                      // 截图
}

func Create(context schema.Context, input CreateParams) (res schema.Response) {
	var (
		err          error
		data         report_schema.Report
		tx           *gorm.DB
		isValidInput bool
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

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
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

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.ErrInvalidParams
		return
	}

	if report_model.IsValidReportType(input.Type) == false {
		err = exception.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	reportInfo := report_model.Report{
		Uid:         context.Uid,
		Title:       input.Title,
		Content:     input.Content,
		Type:        input.Type,
		Screenshots: input.Screenshots,
	}

	if err = tx.Create(&reportInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(reportInfo, &data.ReportPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = reportInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = reportInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func CreateRouter(ctx *gin.Context) {
	var (
		input CreateParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = Create(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
