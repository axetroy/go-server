// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package banner

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
	Image       *string               `json:"image" validate:"omitempty,url,max=255" comment:"图片地址"`                  // 图片 URL
	Href        *string               `json:"href" validate:"omitempty,url,max=255" comment:"图片跳转的地址"`                // 图片跳转的 URL
	Platform    *model.BannerPlatform `json:"platform" validate:"omitempty,max=32,oneof=web app" comment:"平台"`        // 用于哪个平台, web/app
	Description *string               `json:"description" validate:"omitempty,max=255" comment:"描述"`                  // Banner 描述
	Priority    *int                  `json:"priority" validate:"omitempty,gt=0" comment:"优先级"`                       // 优先级，用于排序
	Identifier  *string               `json:"identifier" validate:"omitempty,max=32" comment:"APP 标识符"`               // APP 跳转标识符
	FallbackUrl *string               `json:"fallback_url" validate:"omitempty,url,max=255" comment:"APP 跳转标识符的备选方案"` // APP 跳转标识符的备选方案
}

func Update(c helper.Context, bannerId string, input UpdateParams) (res schema.Response) {
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

		helper.Response(&res, data, nil, err)
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

var UpdateRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateParams
	)

	id := c.Param("banner_id")

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Update(helper.NewContext(&c), id, input)
	})
})
