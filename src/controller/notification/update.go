package notification

import (
	"errors"
	"github.com/asaskevich/govalidator"
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

type UpdateParams struct {
	Tittle  *string `json:"tittle"`  // 公告标题
	Content *string `json:"content"` // 公告内容
	Note    *string `json:"note"`    // 备注
}

func Update(context controller.Context, notificationId string, input UpdateParams) (res schema.Response) {
	var (
		err          error
		data         schema.Notification
		tx           *gorm.DB
		isValidInput bool
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	tx = service.Db.Begin()

	adminInfo := model.Admin{
		Id: context.Uid,
	}

	if err = tx.Where(&adminInfo).First(&adminInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	notificationInfo := model.Notification{
		Id: notificationId,
	}

	if err = tx.Where(&notificationInfo).Last(&notificationInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	updateMap := map[string]interface{}{}

	if input.Tittle != nil && len(*input.Tittle) != 0 {
		updateMap["tittle"] = *input.Tittle
	}

	if input.Content != nil && len(*input.Content) != 0 {
		updateMap["content"] = *input.Content
	}

	if input.Note != nil {
		updateMap["note"] = *input.Note
	}

	if err = tx.Model(&notificationInfo).Updates(updateMap).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(notificationInfo, &data.NotificationPure); err != nil {
		return
	}

	data.CreatedAt = notificationInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = notificationInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateRouter(context *gin.Context) {
	var (
		input UpdateParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = Update(controller.Context{
		Uid: context.GetString("uid"),
	}, context.Param("id"), input)
}
