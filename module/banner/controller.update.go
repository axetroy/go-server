// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/banner/banner_model"
	"github.com/axetroy/go-server/module/banner/banner_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type UpdateParams struct {
	Image       *string                      `json:"image"`        // 图片 URL
	Href        *string                      `json:"href"`         // 图片跳转的 URL
	Platform    *banner_model.BannerPlatform `json:"platform"`     // 用于哪个平台, web/app
	Description *string                      `json:"description"`  // Banner 描述
	Priority    *int                         `json:"priority"`     // 优先级，用于排序
	Identifier  *string                      `json:"identifier"`   // APP 跳转标识符
	FallbackUrl *string                      `json:"fallback_url"` // APP 跳转标识符的备选方案
}

func Update(context schema.Context, bannerId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         banner_schema.Banner
		tx           *gorm.DB
		shouldUpdate bool
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
				err = common_error.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = common_error.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	adminInfo := admin_model.Admin{
		Id: context.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	if !adminInfo.IsSuper {
		err = admin.ErrAdminNotSuper
		return
	}

	bannerInfo := banner_model.Banner{
		Id: bannerId,
	}

	if err = tx.First(&bannerInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrBannerNotExist
			return
		}
		return
	}

	updateModel := banner_model.Banner{}

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
		if err = tx.Model(&bannerInfo).UpdateColumns(&updateModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = ErrBannerNotExist
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

func UpdateRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	id := ctx.Param("banner_id")

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	res = Update(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id, input)
}
