// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/address/address_model"
	"github.com/axetroy/go-server/module/address/address_schema"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
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

func Update(context schema.Context, addressId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         address_schema.Address
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
				err = exception.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil || !shouldUpdate {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	userInfo := user_model.User{
		Id: context.Uid,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	addressInfo := address_model.Address{
		Id:  addressId,
		Uid: context.Uid,
	}

	if err = tx.First(&addressInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrAddressNotExist
			return
		}
		return
	}

	updateModel := address_model.Address{}

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
			err = ErrAddressInvalidProvinceCode
			return
		}

		shouldUpdate = true
		updateModel.ProvinceCode = *input.ProvinceCode

	}

	if input.CityCode != nil {
		// 校验 城市代码
		if _, ok := CityCode[*input.CityCode]; !ok {
			err = ErrAddressInvalidCityCode
			return
		}

		shouldUpdate = true
		updateModel.CityCode = *input.CityCode
	}

	if input.AreaCode != nil {
		// 校验 区域代码
		if _, ok := CountryCode[*input.AreaCode]; !ok {
			err = ErrAddressInvalidAreaCode
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
		if err = tx.Model(&addressInfo).UpdateColumns(&updateModel).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = ErrAddressNotExist
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

func UpdateRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	id := ctx.Param("address_id")

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = exception.ErrInvalidParams
		return
	}

	res = Update(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id, input)
}
