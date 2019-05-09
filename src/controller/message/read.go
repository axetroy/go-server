package message

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

func MarkRead(context controller.Context, id string) (res schema.Response) {
	var (
		err  error
		data schema.Message
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = service.Db.Begin()

	MessageInfo := model.Message{
		Id:  id,
		Uid: context.Uid,
	}

	if err = tx.Where(&MessageInfo).Last(&MessageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if er := mapstructure.Decode(MessageInfo, &data.MessagePure); er != nil {
		err = er
		return
	}

	data.CreatedAt = MessageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = MessageInfo.UpdatedAt.Format(time.RFC3339Nano)

	now := time.Now()

	if err = tx.Model(&MessageInfo).UpdateColumn(model.Message{
		Read:   true,
		ReadAt: &now,
	}).Error; err != nil {
		return
	}

	return
}

func ReadRouter(context *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	id := context.Param("id")

	res = MarkRead(controller.Context{
		Uid: context.GetString("uid"),
	}, id)
}
