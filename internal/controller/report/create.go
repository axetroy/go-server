// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package report

import (
	"errors"
	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateParams struct {
	Title       string           `json:"title" valid:"required~请填写标题"`     // 标题
	Content     string           `json:"content" valid:"required~请填写反馈内容"` // 内容
	Type        model.ReportType `json:"type" valid:"required~请填写反馈类型"`    // 反馈类型
	Screenshots []string         `json:"screenshots"`                      // 截图
}

func Create(c controller.Context, input CreateParams) (res schema.Response) {
	var (
		err  error
		data schema.Report
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if model.IsValidReportType(input.Type) == false {
		err = exception.InvalidParams
		return
	}

	tx = database.Db.Begin()

	reportInfo := model.Report{
		Uid:         c.Uid,
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

func CreateRouter(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Create(controller.NewContext(c), input)
}
