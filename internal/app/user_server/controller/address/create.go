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

type CreateAddressParams struct {
	Name         string  `json:"name" valid:"required~请填写收货人"`           // 收货人
	Phone        string  `json:"phone" valid:"required~请输入收货人电话号码"`      // 收货人手机号
	ProvinceCode string  `json:"province_code" valid:"required~请选择省份"`   // 省份代码
	CityCode     string  `json:"city_code" valid:"required~请选择城市"`       // 城市代码
	AreaCode     string  `json:"area_code" valid:"required~请选择区域"`       // 区域代码
	StreetCode   string  `json:"street_code" valid:"required~请选择街道/乡/镇"` // 街道/乡/镇
	Address      string  `json:"address" valid:"required~请输入详细地址"`       // 详细的地址
	IsDefault    *bool   `json:"is_default"`                             // 是否是默认地址
	Note         *string `json:"note"`                                   // 备注/标签
}

func Create(c helper.Context, input CreateAddressParams) (res schema.Response) {
	var (
		err       error
		data      schema.Address
		tx        *gorm.DB
		isDefault = false
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
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if area.IsValid(input.ProvinceCode, input.CityCode, input.AreaCode, input.StreetCode) == false {
		err = exception.InvalidParams.New("无效的城市码")
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

	if input.IsDefault != nil {
		isDefault = *input.IsDefault

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
	} else {
		firstAddress := model.Address{
			Uid: c.Uid,
		}
		if err = tx.Where(&firstAddress).First(&firstAddress).Error; err != nil {
			// 如果还没有设置过地址，那么这次设置就是默认地址
			if err == gorm.ErrRecordNotFound {
				err = nil
				isDefault = true
			} else {
				return
			}
		}
	}

	AddressInfo := model.Address{
		Uid:          c.Uid,
		Name:         input.Name,
		Phone:        input.Phone,
		ProvinceCode: input.ProvinceCode,
		CityCode:     input.CityCode,
		AreaCode:     input.AreaCode,
		StreetCode:   input.StreetCode,
		Address:      input.Address,
		IsDefault:    isDefault,
		Note:         input.Note,
	}

	if err = tx.Create(&AddressInfo).Error; err != nil {
		return
	}

	if er := mapstructure.Decode(AddressInfo, &data.AddressPure); er != nil {
		err = er
		return
	}

	data.CreatedAt = AddressInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = AddressInfo.UpdatedAt.Format(time.RFC3339Nano)
	return
}

var CreateRouter = router.Handler(func(c router.Context) {
	var (
		input CreateAddressParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Create(helper.NewContext(&c), input)
	})
})
