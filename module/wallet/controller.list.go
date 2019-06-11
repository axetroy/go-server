// Copyright 2019 Axetroy. All rights reserved. MIT license.
package wallet

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/module/wallet/wallet_model"
	"github.com/axetroy/go-server/module/wallet/wallet_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetWallets(context schema.Context) (res schema.Response) {
	var (
		err  error
		data []wallet_schema.Wallet
		list []wallet_model.Wallet
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

	// 获取用户信息
	userInfo := user_model.User{Id: context.Uid}

	tx = database.Db.Begin()

	if err = tx.Where(userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	sql := GenerateWalletSQL(QueryParams{
		Id: &userInfo.Id,
	}, 100, false)

	if err = tx.Raw(sql).Scan(&list).Error; err != nil {
		return
	}

	for _, v := range list {
		wallet := wallet_schema.Wallet{}
		mapToSchema(v, &wallet)
		data = append(data, wallet)
	}

	return
}

func GetWalletsRouter(ctx *gin.Context) {
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

	res = GetWallets(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	})
}
