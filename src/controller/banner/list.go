// Copyright 2019 Axetroy. All rights reserved. MIT license.
package banner

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/middleware"
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
	Platform *model.BannerPlatform `json:"platform"` // 根据平台筛选
	Active   *bool                 `json:"active"`   // 是否激活
}

func GetList(context controller.Context, input Query) (res schema.List) {
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

	query := input.Query

	query.Normalize()

	list := make([]model.Banner, 0)

	var total int64

	filter := model.Banner{}

	if input.Platform != nil {
		filter.Platform = *input.Platform
	}

	if input.Active != nil {
		filter.Active = *input.Active
	} else {
		filter.Active = true
	}

	if err = database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Where(&filter).Find(&list).Count(&total).Error; err != nil {
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

	return
}

func GetListRouter(context *gin.Context) {
	var (
		err   error
		res   = schema.List{}
		input Query
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = GetList(controller.Context{
		Uid: context.GetString(middleware.ContextUidField),
	}, input)
}
