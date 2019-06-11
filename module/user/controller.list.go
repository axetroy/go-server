// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
}

func GetList(context schema.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]user_schema.Profile, 0)
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

	list := make([]user_model.User, 0)

	var total int64

	if err = database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Order(query.Sort).Find(&list).Count(&total).Error; err != nil {
		if err.Error() == common_error.ErrEmptyList.Error() {
			err = nil
		} else {
			return
		}
	}

	for _, v := range list {
		d := user_schema.Profile{}
		if er := mapstructure.Decode(v, &d.ProfilePure); er != nil {
			err = er
			return
		}
		d.PayPassword = v.PayPassword != nil
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
