// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package message

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/middleware"
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
	Status *model.MessageStatus `json:"status" form:"status"`
	Read   *bool                `json:"read" form:"read"`
}

type QueryAdmin struct {
	Query
	Uid *string `json:"uid" form:"uid"` // 指定某个用户ID
}

// 用户获取自己的消息列表
func GetMessageListByUser(c helper.Context, input Query) (res schema.Response) {
	var (
		err  error
		data = make([]schema.Message, 0) // 接口输出的数据
		list = make([]model.Message, 0)  // 数据库查询出的原始数据
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

	query := input.Query

	query.Normalize()

	var total int64

	filter := map[string]interface{}{}

	filter["uid"] = c.Uid

	if input.Read != nil {
		filter["read"] = *input.Read
	}

	if input.Status != nil {
		filter["status"] = *input.Status
	}

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Where(filter).Find(&list).Error; err != nil {
		return
	}

	if err = database.Db.Model(model.Message{}).Where(filter).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.Message{}
		if er := mapstructure.Decode(v, &d.MessagePure); er != nil {
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

// 用户获取自己的消息列表
func GetMessageListByAdmin(c helper.Context, input QueryAdmin) (res schema.Response) {
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

	query := input.Query

	query.Normalize()

	filter := model.Message{}

	if input.Uid != nil {
		filter.Uid = *input.Uid
	}

	if input.Read != nil {
		filter.Read = *input.Read
	}

	if input.Status != nil {
		filter.Status = *input.Status
	}

	var total int64

	if err = query.Order(database.Db.Limit(query.Limit).Offset(query.Limit * query.Page)).Where(&filter).Find(&list).Error; err != nil {
		return
	}

	if err = database.Db.Model(&filter).Where(&filter).Count(&total).Error; err != nil {
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

func GetMessageListByUserRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input Query
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = GetMessageListByUser(helper.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, input)
}

var GetMessageListByAdminRouter = router.Handler(func(c router.Context) {
	var (
		input QueryAdmin
	)

	c.ResponseFunc(c.ShouldBindQuery(&input), func() schema.Response {
		return GetMessageListByAdmin(helper.NewContext(&c), input)
	})
})
