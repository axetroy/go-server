// Copyright 2019 Axetroy. All rights reserved. MIT license.
package menu

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func DeleteMenuById(id string) {
	b := model.Banner{}
	database.DeleteRowByTable(b.TableName(), "id", id)
}

func Delete(context controller.Context, menuId string) (res schema.Response) {
	var (
		err  error
		data schema.Menu
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

	tx = database.Db.Begin()

	adminInfo := model.Admin{Id: context.Uid}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	menuInfo := model.Menu{
		Id: menuId,
	}

	if err = tx.First(&menuInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AddressNotExist
			return
		}
		return
	}

	if err = tx.Delete(model.Menu{
		Id: menuInfo.Id,
	}).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(menuInfo, &data.MenuPure); err != nil {
		return
	}

	data.CreatedAt = menuInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = menuInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func DeleteRouter(c *gin.Context) {
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

	id := c.Param("menu_id")

	res = Delete(controller.NewContext(c), id)
}
