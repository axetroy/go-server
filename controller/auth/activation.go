package auth

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type ActivationParams struct {
	Code string `valid:"Required;" json:"code"`
}

func Activation(input ActivationParams) (res response.Response) {
	var (
		err error
		tx  *gorm.DB
		uid string // 激活码对应的用户ID
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
		} else {
			res.Status = response.StatusSuccess
		}
	}()

	if uid, err = redis.ActivationCode.Get(input.Code).Result(); err != nil {
		err = exception.InvalidActiveCode
		return
	}

	tx = orm.DB.Begin()

	userInfo := model.User{Id: uid}

	if err = tx.Where(&userInfo).Find(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 如果用户的状态不是未激活的话
	if userInfo.Status != model.UserStatusInactivated {
		err = exception.UserHaveActive
		return
	}

	// 更新激活状态
	tx.Model(&userInfo).Update("status", model.UserStatusInit)

	// delete code from redis
	if err = redis.ActivationCode.Del(input.Code).Err(); err != nil {
		return
	}
	return
}

func ActivationRouter(context *gin.Context) {
	var (
		input ActivationParams
		err   error
		res   = response.Response{}
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

	res = Activation(input)
}
