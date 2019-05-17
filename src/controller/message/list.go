package message

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/gin-gonic/gin"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	Status model.MessageStatus `json:"status" form:"status"` // TODO: 添加条件筛选
}

type QueryAdmin struct {
	Query
	Uid string `json:"uid" form:"uid"` // 指定某个用户ID
}

// 用户获取自己的消息列表
func GetList(context controller.Context, input Query) (res schema.List) {
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

	if err = service.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Where(model.Message{Uid: context.Uid}).Find(&list).Count(&total).Error; err != nil {
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

	return
}

// 用户获取自己的消息列表
func GetListByAdmin(context controller.Context, input QueryAdmin) (res schema.List) {
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

	if err = service.Db.Limit(query.Limit).Offset(query.Limit * query.Page).Find(&list).Count(&total).Error; err != nil {
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

	return
}

func GetListRouter(context *gin.Context) {
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

	res = GetList(controller.Context{
		Uid: context.GetString("uid"),
	}, input)
}

func GetListAdminRouter(context *gin.Context) {
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
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindQuery(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = GetListByAdmin(controller.Context{
		Uid: context.GetString("uid"),
	}, input)
}
