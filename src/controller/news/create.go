package news

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateNewParams struct {
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Type    model.NewsType `json:"type"`
	Tags    []string       `json:"tags"`
}

func Create(context controller.Context, input CreateNewParams) (res schema.Response) {
	var (
		err  error
		data schema.News
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
	if !model.IsValidNewsType(input.Type) {
		err = exception.NewsInvalidType
		return
	}

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

	NewsInfo := model.News{
		Author:  context.Uid,
		Title:   input.Title,
		Content: input.Content,
		Type:    input.Type,
		Tags:    input.Tags,
		Status:  model.NewsStatusActive,
	}

	if err = tx.Create(&NewsInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(NewsInfo, &data.NewsPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = NewsInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = NewsInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

func CreateRouter(context *gin.Context) {
	var (
		input CreateNewParams
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
		Uid: context.GetString(middleware.ContextUidField),
	}, input)
}
