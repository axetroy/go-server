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

func DeleteMessageById(id string) {
	service.DeleteRowByTable("message", "id", id)
}

func DeleteByAdmin(context controller.Context, messageId string) (res schema.Response) {
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
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = service.Db.Begin()

	adminInfo := model.Admin{Id: context.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	messageInfo := model.Message{
		Id: messageId,
	}

	if err = tx.First(&messageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.MessageNotExist
			return
		}
		return
	}

	if err = tx.Delete(model.Message{Id: messageInfo.Id}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(messageInfo, &data.MessagePure); err != nil {
		return
	}

	data.CreatedAt = messageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = messageInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func DeleteByAdminRouter(context *gin.Context) {
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

	res = DeleteByAdmin(controller.Context{
		Uid: context.GetString("uid"),
	}, id)
}

func DeleteByUser(context controller.Context, messageId string) (res schema.Response) {
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
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = service.Db.Begin()

	messageInfo := model.Message{
		Id:  messageId,
		Uid: context.Uid,
	}

	if err = tx.First(&messageInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.MessageNotExist
			return
		}
		return
	}

	if err = tx.Delete(model.Message{Id: messageInfo.Id}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(messageInfo, &data.MessagePure); err != nil {
		return
	}

	data.CreatedAt = messageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = messageInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func DeleteByUserRouter(context *gin.Context) {
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

	res = DeleteByUser(controller.Context{
		Uid: context.GetString("uid"),
	}, id)
}
