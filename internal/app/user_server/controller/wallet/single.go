// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package wallet

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"strings"
)

func IsValidWallet(walletName string) bool {
	for _, validWallet := range model.Wallets {
		// 有效币种验证忽略大小写
		if validWallet == strings.ToUpper(walletName) {
			return true
		}
	}
	return false
}

func GetWallet(c helper.Context, currencyName string) (res schema.Response) {
	var (
		err  error
		data schema.Wallet
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

	// 获取用户信息
	userInfo := model.User{Id: c.Uid}

	tx = database.Db.Begin()

	if err = tx.Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	walletInfo := model.Wallet{
		Id: userInfo.Id,
	}

	// 检查是否是有效的钱包
	if !IsValidWallet(strings.ToUpper(currencyName)) {
		err = exception.InvalidWallet
		return
	}

	if err = tx.Table("wallet_"+strings.ToLower(currencyName)).Where("id = ?", c.Uid).Scan(&walletInfo).Error; err != nil {
		return
	}

	mapToSchema(walletInfo, &data)

	return
}

var GetWalletRouter = router.Handler(func(c router.Context) {
	currency := c.Param("currency")

	c.ResponseFunc(nil, func() schema.Response {
		return GetWallet(helper.NewContext(&c), currency)
	})
})
