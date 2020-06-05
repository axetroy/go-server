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

type CreateParams struct {
	Image       string               `json:"image" validate:"required,url,max=255" comment:"图片地址"`                   // 图片 URL
	Href        string               `json:"href" validate:"required,url,max=255" comment:"图片跳转的地址"`                 // 图片跳转的 URL
	Platform    model.BannerPlatform `json:"platform" validate:"required,max=32,oneof=pc app" comment:"平台"`          // 用于哪个平台, web/app
	Description *string              `json:"description" validate:"omitempty,max=255" comment:"描述"`                  // Banner 描述
	Priority    *int                 `json:"priority" validate:"omitempty,gt=0" comment:"优先级"`                       // 优先级，用于排序
	Identifier  *string              `json:"identifier" validate:"omitempty,max=32" comment:"APP 标识符"`               // APP 跳转标识符
	FallbackUrl *string              `json:"fallback_url" validate:"omitempty,url,max=255" comment:"APP 跳转标识符的备选方案"` // APP 跳转标识符的备选方案
}

func Create(c helper.Context, input CreateParams) (res schema.Response) {
	var (
		err  error
		data schema.Banner
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

	if input.Platform == model.BannerPlatformPc {
		// PC 端
	} else if input.Platform == model.BannerPlatformApp {
		// 移动端
	} else {
		err = exception.BannerInvalidPlatform
		return
	}

	bannerInfo := model.Banner{
		// require
		Image:    input.Image,
		Href:     input.Href,
		Platform: input.Platform,
		// optional
		Description: input.Description,
		Priority:    input.Priority,
		Identifier:  input.Identifier,
		FallbackUrl: input.FallbackUrl,
	}

	if err = tx.Create(&bannerInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(bannerInfo, &data.BannerPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = bannerInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = bannerInfo.UpdatedAt.Format(time.RFC3339Nano)

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
