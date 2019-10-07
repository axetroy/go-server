// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
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

func Update(context controller.Context, addressId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.Address
		tx           *gorm.DB
		shouldUpdate bool
		isValidInput bool
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	tx = database.Db.Begin()

	userInfo := model.User{
		Id: context.Uid,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	addressInfo := model.Address{
		Id:  addressId,
		Uid: context.Uid,
	}

	if err = tx.First(&addressInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AddressNotExist
			return
		}
		return
	}

	updateModel := model.Address{}

	if input.Name != nil {
		shouldUpdate = true
		updateModel.Name = *input.Name
	}

	if input.Phone != nil {
		shouldUpdate = true
		updateModel.Phone = *input.Phone
	}

	if input.ProvinceCode != nil {
		// 校验 省份代码
		if _, ok := ProvinceCode[*input.ProvinceCode]; !ok {
			err = exception.AddressInvalidProvinceCode
			return
		}

		shouldUpdate = true
		updateModel.ProvinceCode = *input.ProvinceCode

	}

	if input.CityCode != nil {
		// 校验 城市代码
		if _, ok := CityCode[*input.CityCode]; !ok {
			err = exception.AddressInvalidCityCode
			return
		}

		shouldUpdate = true
		updateModel.CityCode = *input.CityCode
	}

	if input.AreaCode != nil {
		// 校验 区域代码
		if _, ok := CountryCode[*input.AreaCode]; !ok {
			err = exception.AddressInvalidAreaCode
			return
		}

		shouldUpdate = true
		updateModel.AreaCode = *input.AreaCode
	}

	if input.IsDefault != nil {
		shouldUpdate = true
		updateModel.IsDefault = *input.IsDefault
	}

	if shouldUpdate {
		if err = tx.Model(&addressInfo).Updates(&updateModel).Error; err != nil {
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

func UpdateRouter(c *gin.Context) {
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
		c.JSON(http.StatusOK, res)
	}()

	id := c.Param("address_id")

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Update(controller.NewContext(c), id, input)
}
