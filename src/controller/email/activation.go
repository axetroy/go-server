// Copyright 2019 Axetroy. All rights reserved. MIT license.
package email

import (
	"errors"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/service/email"
	"github.com/axetroy/go-server/src/service/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

type SendActivationEmailParams struct {
	To string `json:"to"` // 发送给谁
}

func GenerateActivationCode(uid string) string {
	// 生成重置码
	activationCode := "activation-" + uid
	return activationCode
}

func SendActivationEmail(input SendActivationEmailParams) (res schema.Response) {
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Status = schema.StatusSuccess
		}
	}()

	userInfo := model.User{
		Email: &input.To,
	}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if userInfo.Status != model.UserStatusInactivated {
		err = exception.UserHaveActive
		return
	}

	// generate activation code
	activationCode := GenerateActivationCode(userInfo.Id)

	// set activationCode to redis
	if err = redis.ActivationCodeClient.Set(activationCode, userInfo.Id, time.Minute*30).Err(); err != nil {
		return
	}

	e := email.NewMailer()

	// send email
	if err = e.SendActivationEmail(input.To, activationCode); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ActivationCodeClient.Del(activationCode).Err()
		return
	}

	return
}

func SendActivationEmailRouter(context *gin.Context) {
	var (
		input SendActivationEmailParams
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

	res = SendActivationEmail(input)
}
