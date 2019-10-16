// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
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
	Image       *string               `json:"image"`        // 图片 URL
	Href        *string               `json:"href"`         // 图片跳转的 URL
	Platform    *model.BannerPlatform `json:"platform"`     // 用于哪个平台, web/app
	Description *string               `json:"description"`  // Banner 描述
	Priority    *int                  `json:"priority"`     // 优先级，用于排序
	Identifier  *string               `json:"identifier"`   // APP 跳转标识符
	FallbackUrl *string               `json:"fallback_url"` // APP 跳转标识符的备选方案
}

func Update(c controller.Context, bannerId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.Banner
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

	adminInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	if !adminInfo.IsSuper {
		err = exception.AdminNotSuper
		return
	}

	bannerInfo := model.Banner{
		Id: bannerId,
	}

	if err = tx.First(&bannerInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.BannerNotExist
			return
		}
		return
	}

	updateModel := model.Banner{}

	if input.Image != nil {
		shouldUpdate = true
		updateModel.Image = *input.Image
	}

	if input.Href != nil {
		shouldUpdate = true
		updateModel.Href = *input.Href
	}

	if input.Platform != nil {
		shouldUpdate = true
		updateModel.Platform = *input.Platform

	}

	if input.Description != nil {
		shouldUpdate = true
		updateModel.Description = input.Description
	}

	if input.Priority != nil {
		shouldUpdate = true
		updateModel.Priority = input.Priority
	}

	if input.Identifier != nil {
		shouldUpdate = true
		updateModel.Identifier = input.Identifier
	}

	if input.FallbackUrl != nil {
		shouldUpdate = true
		updateModel.FallbackUrl = input.FallbackUrl
	}

	if shouldUpdate {
		if err = tx.Model(&bannerInfo).Updates(&updateModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.BannerNotExist
				return
			}
			return
		}
	}

	if err = mapstructure.Decode(bannerInfo, &data.BannerPure); err != nil {
		return
	}

	data.CreatedAt = bannerInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = bannerInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	id := c.Param("banner_id")

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Update(controller.NewContext(c), id, input)
}
