// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/module/news/news_model"
	"github.com/axetroy/go-server/module/news/news_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	Status *news_model.NewsStatus `json:"status" form:"status"`
	Type   *news_model.NewsType   `json:"type" form:"type"`
}

func GetList(input Query) (res schema.List) {
	var (
		err  error
		data = make([]news_schema.News, 0) // 接口输出的数据
		list = make([]news_model.News, 0)  // 数据库查询返回的原始数据
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
				err = exception.ErrUnknown
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

	var total int64

	filter := news_model.News{}

	if input.Status != nil {
		filter.Status = *input.Status
	}

	if input.Type != nil {
		filter.Type = *input.Type
	}

	if err = database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Where(&filter).Find(&list).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := news_schema.News{}
		if er := mapstructure.Decode(v, &d.NewsPure); er != nil {
			err = er
			return
		}
		d.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		d.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		data = append(data, d)
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit

	return
}

func GetListRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindQuery(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = GetList(input)
}
