package auth

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/password"
	"github.com/axetroy/go-server/services/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type ResetPasswordParams struct {
	Code        string `json:"code"`
	NewPassword string `json:"new_password"`
}

func ResetPassword(input ResetPasswordParams) (res response.Response) {
	var (
		err error
		tx  *gorm.DB
		uid string // 重置码对应的uid
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
			res.Data = false
		} else {
			res.Status = response.StatusSuccess
			res.Data = true
		}
	}()

	if uid, err = redis.ResetCode.Get(input.Code).Result(); err != nil {
		err = exception.InvalidResetCode
		return
	}

	tx = orm.DB.Begin()

	userInfo := model.User{Id: uid}

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 更新密码
	tx.Model(&userInfo).Update("password", password.Generate(input.NewPassword))

	// delete reset code from redis
	if err = redis.ResetCode.Del(input.Code).Err(); err != nil {
		return
	}

	// TODO: 安全起见，发送一封邮件/短信告知用户
	return
}

func ResetPasswordRouter(context *gin.Context) {
	var (
		input ResetPasswordParams
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

	res = ResetPassword(input)
}
