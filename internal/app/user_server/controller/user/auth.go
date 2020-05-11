package user

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/captcha"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/email"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/service/telephone"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
	"time"
)

func SendAuthEmail(c helper.Context) (res schema.Response) {
	var (
		err  error
		data schema.Profile
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

	userInfo := model.User{}

	if err = tx.Where(&userInfo).Find(&userInfo).Error; err != nil {
		return
	}

	// 如果用户没绑定邮箱，则也没法发送验证码
	if userInfo.Email == nil {
		err = exception.NoData
		return
	}

	// 生成验证码
	activationCode := captcha.GenerateEmailCaptcha()

	// 缓存验证码到 redis
	if err = redis.ClientAuthEmailCode.Set(activationCode, *userInfo.Email, time.Minute*10).Err(); err != nil {
		return
	}

	e := email.NewMailer()

	// send email
	if err = e.SendAuthEmail(*userInfo.Email, activationCode); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ClientAuthEmailCode.Del(activationCode).Err()
		return
	}

	return
}

func SendAuthPhone(c helper.Context) (res schema.Response) {
	var (
		err  error
		data schema.Profile
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

	userInfo := model.User{}

	if err = tx.Where(&userInfo).Find(&userInfo).Error; err != nil {
		return
	}

	// 如果用户没绑定手机号，则也没法发送验证码
	if userInfo.Phone == nil {
		err = exception.NoData
		return
	}

	// 生成验证码
	activationCode := captcha.GeneratePhoneCaptcha()

	// 缓存验证码到 redis
	if err = redis.ClientAuthPhoneCode.Set(activationCode, *userInfo.Phone, time.Minute*10).Err(); err != nil {
		return
	}

	if err = telephone.GetClient().SendAuthCode(*userInfo.Phone, activationCode); err != nil {
		// 如果发送失败，则删除
		_ = redis.ClientAuthPhoneCode.Del(activationCode).Err()
		return
	}

	return
}

func SendAuthEmailRouter(c *gin.Context) {
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

	res = SendAuthEmail(helper.NewContext(c))
}

func SendAuthPhoneRouter(c *gin.Context) {
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

	res = SendAuthPhone(helper.NewContext(c))
}
