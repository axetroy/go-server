// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message

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
	Status *model.MessageStatus `json:"status" url:"status" validate:"omitempty,number" comment:"状态"`
	Read   *bool                `json:"read" url:"read" validate:"omitempty" comment:"是否已读"`
	Uid    *string              `json:"uid" url:"uid" validate:"omitempty" comment:"用户ID"`
}

// 用户获取自己的消息列表
func GetMessageListByAdmin(c helper.Context, query Query) (res schema.Response) {
	var (
		err  error
		data = make([]schema.MessageAdmin, 0) // 接口输出的数据
		list = make([]model.Message, 0)       // 数据库查询出的原始数据
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

	if query.Uid != nil {
		filter["uid"] = *query.Uid
	}

	if query.Read != nil {
		filter["read"] = *query.Read
	}

	if query.Status != nil {
		filter["status"] = *query.Status
	}

	var total int64

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Where(filter).Find(&list).Error; err != nil {
		return
	}

	if err = database.Db.Model(model.Message{}).Where(filter).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.MessageAdmin{}
		if er := mapstructure.Decode(v, &d.MessagePureAdmin); er != nil {
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
	meta.Sort = query.Sort

	return
}

var GetMessageListByAdminRouter = router.Handler(func(c router.Context) {
	var (
		input Query
	)

	c.ResponseFunc(c.ShouldBindQuery(&input), func() schema.Response {
		return GetMessageListByAdmin(helper.NewContext(&c), input)
	})
})
