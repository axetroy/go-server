// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package menu

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

func GetMenu(ctx helper.Context, id string) (res schema.Response) {
	var (
		err  error
		data = schema.Menu{}
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

	menuInfo := model.Menu{
		Id: id,
	}

	if err = database.Db.First(&menuInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	if err = mapstructure.Decode(menuInfo, &data.MenuPure); err != nil {
		return
	}

	if len(data.Accession) == 0 {
		data.Accession = []string{}
	}

	data.Children = []schema.Menu{}

	data.CreatedAt = menuInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = menuInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var GetMenuRouter = router.Handler(func(c router.Context) {
	id := c.Param("menu_id")

	c.ResponseFunc(nil, func() schema.Response {
		return GetMenu(helper.NewContext(&c), id)
	})
})
