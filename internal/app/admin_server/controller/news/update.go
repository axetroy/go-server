// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package news

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type UpdateParams struct {
	Title   *string           `json:"title" validate:"omitempty,max=32" comment:"标题"`
	Content *string           `json:"content" validate:"omitempty" comment:"内容"`
	Type    *model.NewsType   `json:"type" validate:"omitempty,max=32" comment:"类型"`
	Tags    *[]string         `json:"tags" validate:"omitempty" comment:"标题"`
	Status  *model.NewsStatus `json:"status" validate:"omitempty" comment:"状态"`
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

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

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
