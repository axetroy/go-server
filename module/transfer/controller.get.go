// Copyright 2019 Axetroy. All rights reserved. MIT license.
package transfer

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/transfer/transfer_model"
	"github.com/axetroy/go-server/module/transfer/transfer_schema"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func GetDetail(context schema.Context, transferId string) (res schema.Response) {
	var (
		err  error
		tx   *gorm.DB
		data = transfer_schema.TransferLog{}
	)

	defer func() {
		if r := recover(); r != nil {
			switch t := r.(type) {
			case string:
				err = errors.New(t)
			case error:
				err = t
			default:
				err = common_error.ErrUnknown
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
			res.Status = schema.StatusSuccess
			res.Data = data
		}
	}()

	tx = database.Db.Begin()

	userInfo := user_model.User{Id: context.Uid}

	if err = tx.Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
		}
		return
	}

	log := transfer_model.TransferLog{}

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
			err = common_error.ErrNoPermission
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

func GetDetailRouter(ctx *gin.Context) {
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

	res = GetDetail(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, ctx.Param("transfer_id"))
}
