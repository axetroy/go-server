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
	"strings"
)

func IsValidWallet(walletName string) bool {
	for _, validWallet := range wallet_model.Wallets {
		// 有效币种验证忽略大小写
		if validWallet == strings.ToUpper(walletName) {
			return true
		}
	}
	return false
}

func GetWallet(context schema.Context, currencyName string) (res schema.Response) {
	var (
		err  error
		data wallet_schema.Wallet
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

	if err = tx.Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = user_error.ErrUserNotExist
		}
		return
	}

	walletInfo := wallet_model.Wallet{
		Id: userInfo.Id,
	}

	// 检查是否是有效的钱包
	if IsValidWallet(strings.ToUpper(currencyName)) == false {
		err = ErrInvalidWallet
		return
	}

	if err = tx.Table("wallet_"+strings.ToLower(currencyName)).Where("id = ?", context.Uid).Scan(&walletInfo).Error; err != nil {
		return
	}

	mapToSchema(walletInfo, &data)

	return
}

func GetWalletRouter(ctx *gin.Context) {
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

	currency := ctx.Param("currency")

	res = GetWallet(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, currency)
}
