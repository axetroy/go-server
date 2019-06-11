// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/news/news_model"
	"github.com/axetroy/go-server/module/news/news_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateNewParams struct {
	Title   string              `json:"title"`
	Content string              `json:"content"`
	Type    news_model.NewsType `json:"type"`
	Tags    []string            `json:"tags"`
}

func Create(context schema.Context, input CreateNewParams) (res schema.Response) {
	var (
		err  error
		data news_schema.News
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
				err = common_error.ErrUnknown
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
	if !news_model.IsValidNewsType(input.Type) {
		err = ErrNewsInvalidType
		return
	}

	tx = database.Db.Begin()

	adminInfo := admin_model.Admin{
		Id: context.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	if !adminInfo.IsSuper {
		err = admin.ErrAdminNotSuper
		return
	}

	NewsInfo := news_model.News{
		Author:  context.Uid,
		Title:   input.Title,
		Content: input.Content,
		Type:    input.Type,
		Tags:    input.Tags,
		Status:  news_model.NewsStatusActive,
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

func CreateRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	res = Create(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}
