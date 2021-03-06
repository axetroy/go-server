// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package report

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type CreateParams struct {
	Title       string           `json:"title" validate:"required,max=32" comment:"标题"`  // 标题
	Content     string           `json:"content" validate:"required" comment:"内容"`       // 内容
	Type        model.ReportType `json:"type" validate:"required,max=12" comment:"反馈类型"` // 反馈类型
	Screenshots []string         `json:"screenshots" validate:"omitempty" comment:"截图"`  // 截图
}

func Create(c helper.Context, input CreateParams) (res schema.Response) {
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

		helper.Response(&res, data, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
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

	if err = mapstructure.Decode(reportInfo, &data.ReportPure); err != nil {
		return
	}

	data.CreatedAt = reportInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = reportInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var CreateRouter = router.Handler(func(c router.Context) {
	var (
		input CreateParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Create(helper.NewContext(&c), input)
	})
})
