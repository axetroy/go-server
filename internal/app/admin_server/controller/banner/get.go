// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package banner

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

func GetBanner(id string) (res schema.Response) {
	var (
		err  error
		data = schema.Banner{}
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

		helper.Response(&res, data, nil, err)
	}()

	bannerInfo := model.Banner{
		Id: id,
	}

	if err = database.Db.First(&bannerInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.BannerNotExist
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

var GetBannerRouter = router.Handler(func(c router.Context) {
	id := c.Param("banner_id")

	c.ResponseFunc(nil, func() schema.Response {
		return GetBanner(id)
	})
})
