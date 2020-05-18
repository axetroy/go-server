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

type UpdateParams struct {
	Status *model.ReportStatus `json:"status" valid:"required~请选择要标记的状态"`
}

type UpdateByAdminParams struct {
	UpdateParams
	Locked *bool `json:"locked"` // 是否锁定
}

func UpdateByAdmin(c helper.Context, reportId string, input UpdateByAdminParams) (res schema.Response) {
	var (
		err          error
		data         schema.Report
		tx           *gorm.DB
		shouldUpdate bool
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
			if err != nil || !shouldUpdate {
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
		Id: reportId,
	}

	if err = tx.First(&reportInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
			return
		}
		return
	}

	updatedModel := model.Report{}

	if input.Status != nil {
		updatedModel.Status = *input.Status
		shouldUpdate = true
	}

	if input.Locked != nil {
		updatedModel.Locked = *input.Locked
		shouldUpdate = true
	}

	if shouldUpdate == false {
		return
	}

	if err = tx.Model(&reportInfo).Where(&model.Report{
		Id: reportId,
	}).Update(updatedModel).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(reportInfo, &data.ReportPure); err != nil {
		return
	}

	data.CreatedAt = reportInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = reportInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var UpdateByAdminRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateByAdminParams
	)

	reportId := c.Param("report_id")

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return UpdateByAdmin(helper.NewContext(&c), reportId, input)
	})
})
