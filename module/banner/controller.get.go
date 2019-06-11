// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
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

func GetBanner(id string) (res schema.Response) {
	var (
		err  error
		data = banner_schema.Banner{}
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

		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	bannerInfo := banner_model.Banner{
		Id: id,
	}

	if err = database.Db.First(&bannerInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrBannerNotExist
		}
		return
	}

	if err = mapstructure.Decode(bannerInfo, &data.BannerPure); err != nil {
		return
	}

	data.CreatedAt = bannerInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = bannerInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetBannerRouter(ctx *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	id := ctx.Param("banner_id")

	res = GetBanner(id)
}
