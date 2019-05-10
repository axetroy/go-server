package address

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type Query struct {
	schema.Query
	//Status model.NewsStatus `json:"status" form:"status"`
}

func GetList(context controller.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]schema.Address, 0) // 输出到外部的结果
		list = make([]model.Address, 0)  // 数据库查询出来的原始结果
		meta = &schema.Meta{}
		tx   *gorm.DB
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

	tx = service.Db.Begin()

	if err = tx.Limit(query.Limit).Offset(query.Limit*query.Page).Where(model.Address{Uid: context.Uid}, context.Uid).Find(&list).Count(&total).Error; err != nil {
		return
	}

	for _, v := range list {
		d := schema.Address{}
		if er := mapstructure.Decode(v, &d.AddressPure); er != nil {
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
