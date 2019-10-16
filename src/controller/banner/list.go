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
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	Platform *model.BannerPlatform `json:"platform" form:"platform"` // 根据平台筛选
	Active   *bool                 `json:"active" form:"active"`     // 是否激活
}

func GetBannerList(c controller.Context, q Query) (res schema.List) {
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

		helper.ResponseList(&res, data, meta, err)
	}()

	query := q.Query

	query.Normalize()

	list := make([]model.Banner, 0)

	filter := map[string]interface{}{}

	if q.Platform != nil {
		filter["platform"] = *q.Platform
	}

	if q.Active != nil {
		filter["active"] = *q.Active
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

func GetBannerListRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.List{}
		query Query
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindQuery(&query); err != nil {
		return
	}

	res = GetBannerList(controller.NewContext(c), query)
}
