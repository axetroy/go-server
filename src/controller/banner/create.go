// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateParams struct {
	Image       string               `json:"image" valid:"required~请填写图片URL"` // 图片 URL
	Href        string               `json:"href" valid:"required~请填写图片跳转链接"` // 图片跳转的 URL
	Platform    model.BannerPlatform `json:"platform" valid:"required~请选择平台"` // 用于哪个平台, web/app
	Description *string              `json:"description"`                     // Banner 描述
	Priority    *int                 `json:"priority"`                        // 优先级，用于排序
	Identifier  *string              `json:"identifier"`                      // APP 跳转标识符
	FallbackUrl *string              `json:"fallback_url"`                    // APP 跳转标识符的备选方案
}

func Create(context controller.Context, input CreateParams) (res schema.Response) {
	var (
		err          error
		data         schema.Banner
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
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		err = exception.WrapValidatorError(err)
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Id: context.Uid,
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

	res = Create(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, input)
}
