// Copyright 2019 Axetroy. All rights reserved. MIT license.
package login

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	Uid     *string `json:"uid" form:"uid"`         // 根据用户 ID 筛选
	Type    *int    `json:"type" form:"type"`       // 根据类型筛选
	Command *int    `json:"command" form:"command"` // 根据登陆命令筛选
	Ip      *string `json:"ip"`                     // 根据 IP 筛选
}

func GetLoginLogs(context controller.Context, q Query) (res schema.List) {
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

		helper.ResponseList(&res, data, meta, err)
	}()

	query := q.Query

	query.Normalize()

	list := make([]model.LoginLog, 0)

	filter := map[string]interface{}{}

	if q.Uid != nil {
		filter["uid"] = *q.Uid
	}

	if q.Type != nil {
		filter["type"] = *q.Type
	}

	if q.Command != nil {
		filter["command"] = *q.Command
	}

	if q.Ip != nil {
		filter["last_ip"] = *q.Ip
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

	return
}

func GetLoginLogsRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.List{}
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

	res = GetLoginLogs(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, query)
}
