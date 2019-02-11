package notification

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

// Query params
type Query struct {
	schema.Query
}

// GetList get notification list
func GetList(context controller.Context, input Query) (res schema.List) {
	var (
		err  error
		data = make([]schema.Notification, 0)
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

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
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

	tx = service.Db.Begin()

	var total int64

	list := make([]model.Notification, 0)

	if err = tx.Table(new(model.Notification).TableName()).Limit(query.Limit).Offset(query.Limit * query.Page).Find(&list).Count(&total).Error; err != nil {
		return
	}

	data = make([]schema.Notification, len(list))

	// TODO: 优化这一块实现
	for index, v := range list {
		current := data[index]
		if er := mapstructure.Decode(v, &current.NotificationPure); er != nil {
			err = er
			return
		}
		current.CreatedAt = v.CreatedAt.Format(time.RFC3339Nano)
		current.UpdatedAt = v.UpdatedAt.Format(time.RFC3339Nano)
		if len(context.Uid) != 0 {

			mark := model.NotificationMark{
				Id:  v.Id,
				Uid: context.Uid,
			}
			if err = tx.Where(&mark).Last(&mark).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					current.Read = false
					current.ReadAt = ""
				}
			} else {
				current.Read = true
				current.ReadAt = mark.CreatedAt.Format(time.RFC3339Nano)
			}
		}
	}

	meta.Total = total
	meta.Num = len(data)
	meta.Page = query.Page
	meta.Limit = query.Limit

	return
}

// GetListRouter get list router
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
