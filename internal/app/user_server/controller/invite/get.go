// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package invite

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func Get(c helper.Context, id string) (res schema.Response) {
	var (
		err  error
		data schema.Invite
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

	inviteDetail := model.InviteHistory{
		Id: id,
	}

	tx = database.Db.Begin()

	// 只能获取跟自己相关的
	if err = tx.Where(model.InviteHistory{Inviter: c.Uid}).Or(model.InviteHistory{Invitee: c.Uid}).First(&inviteDetail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InviteNotExist
			return
		}
	}

	if err = mapstructure.Decode(inviteDetail, &data.InvitePure); err != nil {
		return
	}

	data.CreatedAt = inviteDetail.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = inviteDetail.UpdatedAt.Format(time.RFC3339Nano)

	return
}

// 内部使用, 不对外提供
func GetByStruct(m *model.InviteHistory) (res schema.Response) {
	var (
		err  error
		data schema.Invite
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

	tx = database.Db.Begin()

	if err = tx.Where(m).Last(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InviteNotExist
		}
		return
	}

	if err = mapstructure.Decode(m, &data.InvitePure); err != nil {
		return
	}

	data.CreatedAt = m.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = m.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetRouter(c *gin.Context) {
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

	inviteId := c.Param("invite_id")

	res = Get(helper.NewContext(c), inviteId)
}
