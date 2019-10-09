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
	"github.com/axetroy/go-server/src/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type UpdateParams struct {
	Status *model.ReportStatus `json:"status" valid:"required~请选择要标记的状态"`
}

type UpdateByAdminParams struct {
	UpdateParams
	Locked *bool `json:"locked"` // 是否锁定
}

func Update(context controller.Context, reportId string, input UpdateParams) (res schema.Response) {
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	reportInfo := model.Report{
		Id:  reportId,
		Uid: context.Uid,
	}

	if err = tx.First(&reportInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
			return
		}
		return
	}

	// 如果已被锁定，则无法更新状态
	if reportInfo.Locked {
		err = errors.New("该反馈已被锁定, 无法更新")
		return
	}

	updatedModel := model.Report{}

	if input.Status != nil {
		// 状态不能重复改变, 忽略本次操作.
		if reportInfo.Status == *input.Status {
			return
		}
		updatedModel.Status = *input.Status
		shouldUpdate = true
	}

	if shouldUpdate == false {
		return
	}

	if err = tx.Model(&reportInfo).Where(&model.Report{
		Id:  reportId,
		Uid: context.Uid,
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

func UpdateRouter(c *gin.Context) {
	var (
		input UpdateParams
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

	reportId := c.Param("report_id")

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Update(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, reportId, input)
}

func UpdateByAdmin(context controller.Context, reportId string, input UpdateByAdminParams) (res schema.Response) {
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

		helper.Response(&res, data, err)
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

func UpdateByAdminRouter(c *gin.Context) {
	var (
		input UpdateByAdminParams
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

	reportId := c.Param("report_id")

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = UpdateByAdmin(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, reportId, input)
}
