package invite

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

func Get(context controller.Context, id string) (res schema.Response) {
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

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	inviteDetail := model.InviteHistory{
		Id: id,
	}

	tx = service.Db.Begin()

	if err = tx.Where(&inviteDetail).Last(&inviteDetail).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InviteNotExist
			return
		}
	}

	// 只有跟自己有关的，才能获取详情
	if inviteDetail.Inviter != context.Uid && inviteDetail.Invitee != context.Uid {
		err = exception.NoPermission
		return
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

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = service.Db.Begin()

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

func GetRouter(context *gin.Context) {
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

	inviteId := context.Param("invite_id")

	res = Get(controller.Context{
		Uid: context.GetString("uid"),
	}, inviteId)
}
