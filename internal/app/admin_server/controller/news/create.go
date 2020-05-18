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

type CreateNewParams struct {
	Title   string         `json:"title"`
	Content string         `json:"content"`
	Type    model.NewsType `json:"type"`
	Tags    []string       `json:"tags"`
}

func Create(c helper.Context, input CreateNewParams) (res schema.Response) {
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

		helper.Response(&res, data, nil, err)
	}()

	// 参数校验
	if !model.IsValidNewsType(input.Type) {
		err = exception.NewsInvalidType
		return
	}

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Id: c.Uid,
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
		Author:  c.Uid,
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

var CreateRouter = router.Handler(func(c router.Context) {
	var (
		input CreateNewParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Create(helper.NewContext(&c), input)
	})
})
