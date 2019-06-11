// Copyright 2019 Axetroy. All rights reserved. MIT license.
package address

import (
	"errors"
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

func DeleteAddressById(id string) {
	database.DeleteRowByTable("address", "id", id)
}

func Delete(context schema.Context, addressId string) (res schema.Response) {
	var (
		err  error
		data address_schema.Address
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
				err = exception.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil {
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

	tx = database.Db.Begin()

	userInfo := user_model.User{Id: context.Uid}

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

	if err = tx.Delete(address_model.Address{
		Id: addressInfo.Id,
	}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(addressInfo, &data.AddressPure); err != nil {
		return
	}

	data.CreatedAt = addressInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = addressInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func DeleteRouter(ctx *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	id := ctx.Param("address_id")

	res = Delete(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}
