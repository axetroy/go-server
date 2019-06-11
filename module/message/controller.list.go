// Copyright 2019 Axetroy. All rights reserved. MIT license.
package message

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/message/message_model"
	"github.com/axetroy/go-server/module/message/message_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	Status *message_model.MessageStatus `json:"status" form:"status"`
	Read   *bool                        `json:"read" form:"read"`
}

type QueryAdmin struct {
	Query
	Uid *string `json:"uid" form:"uid"` // 指定某个用户ID
}

// 用户获取自己的消息列表
func GetList(context schema.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]message_schema.Message, 0) // 接口输出的数据
		list = make([]message_model.Message, 0)  // 数据库查询出的原始数据
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

	query := input.Query

	query.Normalize()

	var total int64

	filter := message_model.Message{
		Uid: context.Uid,
	}

	if input.Read != nil {
		filter.Read = *input.Read
	}

	if input.Status != nil {
		filter.Status = *input.Status
	}

	if err = database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Order(query.Sort).Where(&filter).Find(&list).Error; err != nil {
		return
	}

	if err = database.Db.Model(&filter).Where(&filter).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := message_schema.Message{}
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

	return
}

// 用户获取自己的消息列表
func GetListByAdmin(context schema.Context, input QueryAdmin) (res schema.List) {
	var (
		err  error
		data = make([]message_schema.MessageAdmin, 0) // 接口输出的数据
		list = make([]message_model.Message, 0)       // 数据库查询出的原始数据
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

	query := input.Query

	query.Normalize()

	filter := message_model.Message{}

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

	if err = database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Order(query.Sort).Where(&filter).Find(&list).Error; err != nil {
		return
	}

	if err = database.Db.Model(&filter).Where(&filter).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := message_schema.MessageAdmin{}
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
		err = common_error.ErrInvalidParams
		return
	}

	res = GetList(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}

func GetListAdminRouter(ctx *gin.Context) {
	var (
		err   error
		res   = schema.List{}
		input QueryAdmin
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindQuery(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	res = GetListByAdmin(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
