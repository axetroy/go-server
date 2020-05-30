// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package address

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/area"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type UpdateParams struct {
	Name         *string `json:"name"`
	Phone        *string `json:"phone"`
	ProvinceCode *string `json:"province_code"`
	CityCode     *string `json:"city_code"`
	AreaCode     *string `json:"area_code"`
	Address      *string `json:"address"`
	IsDefault    *bool   `json:"is_default"`
}

func Update(c helper.Context, addressId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.Address
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

	userInfo := model.User{
		Id: c.Uid,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	addressInfo := model.Address{
		Id:  addressId,
		Uid: c.Uid,
	}

	if err = tx.First(&addressInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AddressNotExist
			return
		}
		return
	}

	updateModel := map[string]interface{}{}

	if input.Name != nil {
		shouldUpdate = true
		updateModel["name"] = *input.Name
	}

	if input.Phone != nil {
		shouldUpdate = true
		updateModel["phone"] = *input.Phone
	}

	if input.Address != nil {
		shouldUpdate = true
		updateModel["address"] = *input.Address
	}

	if input.ProvinceCode != nil {
		if input.CityCode == nil || input.AreaCode == nil {
			err = exception.InvalidParams
			return
		}

		shouldUpdate = true
		updateModel["province_code"] = *input.ProvinceCode
		updateModel["city_code"] = *input.CityCode
		updateModel["area_code"] = *input.AreaCode

		if area.IsValid(*input.ProvinceCode, *input.CityCode, *input.AreaCode) == false {
			err = exception.InvalidParams
			return
		}
	}

	if input.IsDefault != nil {
		shouldUpdate = true
		updateModel["is_default"] = *input.IsDefault
	}

	if shouldUpdate {
		if err = tx.Model(&addressInfo).Updates(updateModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.AddressNotExist
				return
			}
			return
		}
	}

	if err = mapstructure.Decode(addressInfo, &data.AddressPure); err != nil {
		return
	}

	data.CreatedAt = addressInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = addressInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var UpdateRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateParams
	)

	id := c.Param("address_id")

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Update(helper.NewContext(&c), id, input)
	})
})
