// Copyright 2019 Axetroy. All rights reserved. MIT license.
package news

import (
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/address"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/news/news_model"
	"github.com/axetroy/go-server/module/news/news_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func DeleteNewsById(id string) {
	database.DeleteRowByTable("news", "id", id)
}

func Delete(context schema.Context, addressId string) (res schema.Response) {
	var (
		err  error
		data news_schema.News
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
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	adminInfo := admin_model.Admin{Id: context.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
		}
		return
	}

	newsInfo := news_model.News{
		Id: addressId,
	}

	if err = tx.First(&newsInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = address.ErrAddressNotExist
			return
		}
		return
	}

	if err = tx.Delete(news_model.News{
		Id: newsInfo.Id,
	}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(newsInfo, &data.NewsPure); err != nil {
		return
	}

	data.CreatedAt = newsInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = newsInfo.UpdatedAt.Format(time.RFC3339Nano)

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

	id := ctx.Param("news_id")

	res = Delete(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, id)
}
