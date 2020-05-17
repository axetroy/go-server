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

func GetNews(id string) (res schema.Response) {
	var (
		err  error
		data = schema.News{}
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

		helper.Response(&res, data, nil, err)
	}()

	newsInfo := model.News{
		Id: id,
	}

	if err = database.Db.First(&newsInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NewsNotExist
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

var GetNewsRouter = router.Handler(func(c router.Context) {
	id := c.Param("news_id")

	c.ResponseFunc(nil, func() schema.Response {
		return GetNews(id)
	})
})
