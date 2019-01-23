package notification

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

func MarkRead(context controller.Context, notificationId string) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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
			res.Data = false
			res.Message = err.Error()
		} else {
			res.Data = true
			res.Status = schema.StatusSuccess
		}
	}()

	tx = service.Db.Begin()

	userInfo := model.User{
		Id: context.Uid,
	}

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		// 没有找到用户
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	notificationInfo := model.Notification{
		Id: notificationId,
	}

	// 先获取通知
	if err = tx.Where(&notificationInfo).Last(&notificationInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.NoData
		}
		return
	}

	mark := model.NotificationMark{
		Id:  notificationInfo.Id,
		Uid: context.Uid,
	}

	// 再确认以读表有没有这个用户的已读记录
	if err = tx.Where(&mark).Last(&mark).Error; err != nil {
		// 如果没找到这条记录，则说明没有创建
		// 继续下面的页面
		if err == gorm.ErrRecordNotFound {
			err = nil
		} else {
			return
		}
	} else {
		// 通知已读
		return
	}

	if err = tx.Create(&mark).Error; err != nil {
		return
	}

	return
}

func ReadRouter(context *gin.Context) {
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

	notificationId := context.Param("id")

	res = MarkRead(controller.Context{
		Uid: context.GetString("uid"),
	}, notificationId)
}
