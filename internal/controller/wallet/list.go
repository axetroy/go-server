// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package wallet

import (
	"errors"
	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/helper"
	"github.com/axetroy/go-server/internal/middleware"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func GetWallets(c controller.Context) (res schema.Response) {
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

		helper.Response(&res, data, err)
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

func GetWalletsRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	res = GetWallets(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	})
}
