// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package invite

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
)

type Query struct {
	schema.Query
}

func GetInviteListByUser(query Query) (res schema.Response) {
	var (
		err  error
		data = make([]model.InviteHistory, 0)
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

	filter := map[string]interface{}{}

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Where(filter)).Find(&data).Error; err != nil {
		return
	}

	var total int64

	if err = database.Db.Model(model.InviteHistory{}).Where(filter).Count(&total).Error; err != nil {
		return
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit
	meta.Sort = query.Sort

	return
}

var GetInviteListByUserRouter = router.Handler(func(c router.Context) {
	var (
		input Query
	)

	c.ResponseFunc(c.ShouldBindQuery(&input), func() schema.Response {
		return GetInviteListByUser(input)
	})
})
