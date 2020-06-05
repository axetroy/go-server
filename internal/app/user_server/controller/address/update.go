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
	Name         *string `json:"name" validate:"omitempty,max=32" comment:"收货人"`                    // 收货人
	Phone        *string `json:"phone" validate:"omitempty,numeric,len=11" comment:"电话号码"`          // 收货人手机号
	ProvinceCode *string `json:"province_code" validate:"omitempty,numeric,len=2" comment:"省份代码"`   // 省份代码
	CityCode     *string `json:"city_code" validate:"omitempty,numeric,len=4" comment:"城市代码"`       // 城市代码
	AreaCode     *string `json:"area_code" validate:"omitempty,numeric,len=6" comment:"区域代码"`       // 区域代码
	StreetCode   *string `json:"street_code" validate:"omitempty,numeric,len=9" comment:"街道/乡/镇代码"` // 街道/乡/镇
	Address      *string `json:"address" validate:"omitempty,max=32" comment:"详细地址"`                // 详细的地址
	IsDefault    *bool   `json:"is_default" omitempty:"omitempty"`                                  // 是否是默认地址
	Note         *string `json:"note" validate:"omitempty,max=12"`                                  // 备注/标签
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
		addressInfo.Name = *input.Name
	}

	if input.Phone != nil {
		shouldUpdate = true
		updateModel["phone"] = *input.Phone
		addressInfo.Phone = *input.Phone
	}

	if input.Address != nil {
		shouldUpdate = true
		updateModel["address"] = *input.Address
		addressInfo.Address = *input.Address
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
		addressInfo.ProvinceCode = *input.ProvinceCode
		addressInfo.CityCode = *input.CityCode
		addressInfo.AreaCode = *input.AreaCode

		if area.IsValid(*input.ProvinceCode, *input.CityCode, *input.AreaCode, *input.StreetCode) == false {
			err = exception.InvalidParams.New("无效的城市码")
			return
		}
	}

	if input.IsDefault != nil {
		shouldUpdate = true
		updateModel["is_default"] = *input.IsDefault
		addressInfo.IsDefault = *input.IsDefault

		// 如果要创建一个默认地址
		// 那么就把前面的默认地址修改为false
		if *input.IsDefault == true {
			defaultAddress := model.Address{
				Uid:       c.Uid,
				IsDefault: true,
			}
			if err = tx.Where(&defaultAddress).First(&defaultAddress).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					err = nil
				} else {
					return
				}
			} else {
				// 如果存在了默认地址，则取消它的默认属性
				if err = tx.Model(&defaultAddress).Where("id = ?", defaultAddress.Id).UpdateColumn("is_default", false).Error; err != nil {
					return
				}
			}
		}
	}

	if input.Note != nil {
		shouldUpdate = true
		updateModel["note"] = *input.Note
		addressInfo.Note = input.Note
	}

	if shouldUpdate {
		if err = tx.Model(&addressInfo).Where("id = ?", addressId).Where("uid = ?", c.Uid).Updates(updateModel).Error; err != nil {
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
