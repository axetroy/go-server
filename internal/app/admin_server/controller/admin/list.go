// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin

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
	Status *model.AdminStatus `json:"status" url:"status" validate:"omitempty,number" comment:"管理员状态"` // 管理员状态
}

func GetList(c helper.Context, query Query) (res schema.Response) {
	var (
		err  error
		data = make([]schema.AdminProfile, 0)
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

	list := make([]model.Admin, 0)
	filter := map[string]interface{}{}

	if query.Status != nil {
		filter["status"] = *query.Status
	}

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Where(filter).Find(&list).Error; err != nil {
		return
	}

	var total int64

	if err = database.Db.Model(model.Admin{}).Where(filter).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.AdminProfile{}
		if er := mapstructure.Decode(v, &d.AdminProfilePure); er != nil {
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

var GetListRouter = router.Handler(func(c router.Context) {
	var (
		input Query
	)

	c.ResponseFunc(c.ShouldBindQuery(&input), func() schema.Response {
		return GetList(helper.NewContext(&c), input)
	})
})
