// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package help

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	Status *model.HelpStatus `json:"status" form:"status"` // 根据状态筛选
	Type   *model.HelpType   `json:"type" form:"type"`     // 根据类型筛选
}

func GetHelpList(c helper.Context, q Query) (res schema.Response) {
	var (
		err  error
		data = make([]schema.Help, 0)
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

	query := q.Query

	query.Normalize()

	list := make([]model.Help, 0)

	filter := map[string]interface{}{}

	if q.Status != nil {
		filter["status"] = *q.Status
	}

	if q.Type != nil {
		filter["type"] = *q.Type
	}

	var total int64

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Where(filter).Find(&list).Error; err != nil {
		return
	}

	if err = database.Db.Model(model.Help{}).Where(filter).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.Help{}
		if er := mapstructure.Decode(v, &d.HelpPure); er != nil {
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

func GetHelpListRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
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

	res = GetHelpList(helper.NewContext(c), query)
}
