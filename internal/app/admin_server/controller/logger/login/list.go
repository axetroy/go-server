// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package login

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
	Uid     *string `json:"uid" url:"uid" validate:"omitempty" comment:"用户ID"`                // 根据用户 ID 筛选
	Type    *int    `json:"type" url:"type" validate:"omitempty,number" comment:"类型"`         // 根据类型筛选
	Command *int    `json:"command" url:"command" validate:"omitempty,number" comment:"登陆命令"` // 根据登陆命令筛选
	Ip      *string `json:"ip" url:"ip" validate:"omitempty,ip" comment:"IP"`                 // 根据 IP 筛选
}

func GetLoginLogs(c helper.Context, query Query) (res schema.Response) {
	var (
		err  error
		data = make([]schema.LogLogin, 0)
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

	list := make([]model.LoginLog, 0)

	filter := map[string]interface{}{}

	if query.Uid != nil {
		filter["uid"] = *query.Uid
	}

	if query.Type != nil {
		filter["type"] = *query.Type
	}

	if query.Command != nil {
		filter["command"] = *query.Command
	}

	if query.Ip != nil {
		filter["last_ip"] = *query.Ip
	}

	var total int64

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Where(filter).Find(&list).Error; err != nil {
		return
	}

	if err = database.Db.Model(model.LoginLog{}).Where(filter).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.LogLogin{}
		if er := mapstructure.Decode(v, &d.LogLoginPure); er != nil {
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

var GetLoginLogsRouter = router.Handler(func(c router.Context) {
	var (
		query Query
	)

	c.ResponseFunc(c.ShouldBindQuery(&query), func() schema.Response {
		return GetLoginLogs(helper.NewContext(&c), query)
	})
})
