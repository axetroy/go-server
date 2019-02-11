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

type CreateMessageParams struct {
	Uid     string `json:"uid"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

func Create(context controller.Context, input CreateMessageParams) (res schema.Response) {
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

	// 参数校验

	tx = service.Db.Begin()

	adminInfo := model.Admin{
		Id: context.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	if !adminInfo.IsSuper {
		err = exception.AdminNotSuper
		return
	}

	MessageInfo := model.Message{
		Uid:     input.Uid,
		Tittle:  input.Title,
		Content: input.Content,
		Status:  model.MessageStatusActive,
	}

	if err = tx.Create(&MessageInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(MessageInfo, &data.MessagePure); er != nil {
		err = er
		return
	}

	data.CreatedAt = MessageInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = MessageInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

func CreateRouter(context *gin.Context) {
	var (
		input CreateMessageParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Create(controller.Context{
		Uid: context.GetString("uid"),
	}, input)
}
