// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/banner/banner_model"
	"github.com/axetroy/go-server/module/banner/banner_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	Platform *banner_model.BannerPlatform `json:"platform" form:"platform"` // 根据平台筛选
	Active   *bool                        `json:"active" form:"active"`     // 是否激活
}

func GetList(context schema.Context, q Query) (res schema.List) {
	var (
		err  error
		data = make([]banner_schema.Banner, 0)
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
				err = common_error.ErrUnknown
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
			res.Meta = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
			res.Meta = meta
		}
	}()

	query := q.Query

	query.Normalize()

	list := make([]banner_model.Banner, 0)

	m := map[string]interface{}{}

	if q.Platform != nil {
		m["platform"] = *q.Platform
	}

	if q.Active != nil {
		m["active"] = *q.Active
	} else {
		m["active"] = true
	}

	var total int64

	if err = database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Order(query.Sort).Where(m).Find(&list).Error; err != nil {
		return
	}

	if err = database.Db.Model(banner_model.Banner{}).Where(m).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := banner_schema.Banner{}
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

	return
}

func GetListRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	res = GetList(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, query)
}
