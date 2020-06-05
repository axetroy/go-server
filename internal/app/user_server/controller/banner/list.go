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
	"github.com/mitchellh/mapstructure"
	"time"
)

type Query struct {
	schema.Query
	Platform *model.BannerPlatform `json:"platform" url:"platform" validate:"omitempty,oneof=pc app" comment:"平台"` // 根据平台筛选
	Active   *bool                 `json:"active" url:"active" validate:"omitempty" comment:"是否激活"`                // 是否激活
}

func GetBannerList(c helper.Context, query Query) (res schema.Response) {
	var (
		err  error
		data = make([]schema.Banner, 0)
		meta = &schema.Meta{}
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

		helper.Response(&res, data, meta, err)
	}()

	query.Normalize()

	if err = query.Validate(); err != nil {
		return
	}

	list := make([]model.Banner, 0)

	filter := map[string]interface{}{}

	if query.Platform != nil {
		filter["platform"] = *query.Platform
	}

	if query.Active != nil {
		filter["active"] = *query.Active
	} else {
		filter["active"] = true
	}

	var total int64

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Where(filter).Find(&list).Error; err != nil {
		return
	}

	if err = database.Db.Model(model.Banner{}).Where(filter).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.Banner{}
		if er := mapstructure.Decode(v, &d.BannerPure); er != nil {
			err = er
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(list)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}

var GetBannerListRouter = router.Handler(func(c router.Context) {
	var (
		query Query
	)

	c.ResponseFunc(c.ShouldBindQuery(&query), func() schema.Response {
		return GetBannerList(helper.NewContext(&c), query)
	})
})
