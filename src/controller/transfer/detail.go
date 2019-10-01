// Copyright 2019 Axetroy. All rights reserved. MIT license.
package transfer

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func GetDetail(context controller.Context, transferId string) (res schema.Response) {
	var (
		err  error
		tx   *gorm.DB
		data = schema.TransferLog{}
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

	tx = database.Db.Begin()

	userInfo := model.User{Id: context.Uid}

	if err = tx.Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	log := model.TransferLog{}

	// 联表查询
	// 只能获取自己转给别人的
	sql := GenerateTransferLogSQL(QueryParams{
		Id: &transferId,
	}, 1, false)

	if err = tx.Raw(sql).Scan(&log).Error; err != nil {
		return
	}

	if log.From != context.Uid {
		if log.To != context.Uid {
			// 既不是转账人，也不是收款人, 没有权限获取这条记录
			err = exception.NoPermission
			return
		}
	}

	if err = mapstructure.Decode(log, &data.TransferLogPure); err != nil {
		return
	}

	data.CreatedAt = log.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = log.UpdatedAt.Format(time.RFC3339Nano)
	return
}

func GetDetailRouter(context *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	res = GetDetail(controller.Context{
		Uid: context.GetString(middleware.ContextUidField),
	}, context.Param("transfer_id"))
}
