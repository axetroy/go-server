// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package address

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

func GetDefault(c helper.Context) (res schema.Response) {
	var (
		err  error
		data = schema.Address{}
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

	addressInfo := model.Address{
		Uid:       c.Uid,
		IsDefault: true,
	}

	if err = database.Db.Model(&addressInfo).Where(&addressInfo).First(&addressInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AddressDefaultNotExist
		}
		return
	}

	if err = mapstructure.Decode(addressInfo, &data.AddressPure); err != nil {
		return
	}

	data.CreatedAt = addressInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = addressInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var GetDefaultRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return GetDefault(helper.NewContext(&c))
	})
})
