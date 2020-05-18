// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package news

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type UpdateParams struct {
	Title   *string           `json:"title"`
	Content *string           `json:"content"`
	Type    *model.NewsType   `json:"type"`
	Tags    *[]string         `json:"tags"`
	Status  *model.NewsStatus `json:"status"`
}

func Update(c helper.Context, newsId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.News
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
				err = exception.Unknown
			}
		}

		if tx != nil {
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		helper.Response(&res, data, nil, err)
	}()

	tx = database.Db.Begin()

	adminInfo := model.Admin{Id: c.Uid}

	// 判断管理员是否存在
	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
	}

	newsInfo := model.News{
		Id: newsId,
	}

	if err = tx.First(&newsInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NewsNotExist
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
			err = exception.NewsNotExist
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

var UpdateRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Update(helper.NewContext(&c), c.Param("news_id"), input)
	})
})
