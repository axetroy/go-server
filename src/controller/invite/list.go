// Copyright 2019 Axetroy. All rights reserved. MIT license.
package invite

import (
	"errors"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Query struct {
	schema.Query
}

func GetInviteListByUser(input Query) (res schema.List) {
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

		helper.ResponseList(&res, data, meta, err)
	}()

	query := input.Query

	query.Normalize()

	filter := map[string]interface{}{}

	if err = database.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Where(filter).Order(query.Sort).Find(&data).Error; err != nil {
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

	return
}

func GetInviteListByUserRouter(context *gin.Context) {
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

	res = GetInviteListByUser(input)
}
