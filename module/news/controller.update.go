// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news

import (
	"errors"
	"github.com/axetroy/go-server/exception"
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

type UpdateParams struct {
	Title   *string                `json:"title"`
	Content *string                `json:"content"`
	Type    *news_model.NewsType   `json:"type"`
	Tags    *[]string              `json:"tags"`
	Status  *news_model.NewsStatus `json:"status"`
}

func Update(context schema.Context, newsId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         news_schema.News
		tx           *gorm.DB
		shouldUpdate bool
	)

	defer func() {

		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = exception.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil || !shouldUpdate {
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

	tx = database.Db.Begin()

	adminInfo := admin_model.Admin{Id: context.Uid}

	// 判断管理员是否存在
	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
			return
		}
	}

	newsInfo := news_model.News{
		Id: newsId,
	}

	if err = tx.First(&newsInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrNewsNotExist
			return
		}
		return
	}

	if input.Title != nil {
		shouldUpdate = true
		newsInfo.Title = *input.Title
	}

	if input.Content != nil {
		shouldUpdate = true
		newsInfo.Content = *input.Content
	}

	if input.Type != nil {
		shouldUpdate = true
		newsInfo.Type = *input.Type
	}

	if input.Status != nil {
		shouldUpdate = true
		newsInfo.Status = *input.Status
	}

	if input.Tags != nil {
		shouldUpdate = true
		newsInfo.Tags = *input.Tags
	}

	if err = tx.Save(&newsInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrNewsNotExist
			return
		}
		return
	}

	if err = mapstructure.Decode(newsInfo, &data.NewsPure); err != nil {
		return
	}

	data.CreatedAt = newsInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = newsInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateRouter(ctx *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	id := ctx.Param("news_id")

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = Update(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id, input)
}
