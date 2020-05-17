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
)

func GetWallets(c helper.Context) (res schema.Response) {
	var (
		err  error
		data []schema.Wallet
		list []model.Wallet
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

	if err = tx.Where(userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
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
		wallet := schema.Wallet{}
		mapToSchema(v, &wallet)
		data = append(data, wallet)
	}

	return
}

var GetWalletsRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return GetWallets(helper.NewContext(&c))
	})
})
