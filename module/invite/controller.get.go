// Copyright 2019 Axetroy. All rights reserved. MIT license.
package invite

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/invite/invite_model"
	"github.com/axetroy/go-server/module/invite/invite_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func Get(context schema.Context, id string) (res schema.Response) {
	var (
		err  error
		data invite_schema.Invite
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

	inviteDetail := invite_model.InviteHistory{
		Id: id,
	}

	tx = database.Db.Begin()

	// 只能获取跟自己相关的
	if err = tx.Where(invite_model.InviteHistory{Inviter: context.Uid}).Or(invite_model.InviteHistory{Invitee: context.Uid}).First(&inviteDetail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrInviteNotExist
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
func GetByStruct(m *invite_model.InviteHistory) (res schema.Response) {
	var (
		err  error
		data invite_schema.Invite
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

	tx = database.Db.Begin()

	if err = tx.Where(m).Last(m).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = ErrInviteNotExist
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

func GetRouter(ctx *gin.Context) {
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

	inviteId := ctx.Param("invite_id")

	res = Get(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, inviteId)
}
