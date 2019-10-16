package auth

import (
	"errors"
	"github.com/axetroy/go-server/core/controller"
	"github.com/axetroy/go-server/core/exception"
	"github.com/axetroy/go-server/core/helper"
	"github.com/axetroy/go-server/core/schema"
	"github.com/axetroy/go-server/core/service/email"
	"github.com/axetroy/go-server/core/service/redis"
	"github.com/axetroy/go-server/core/service/telephone"
	"github.com/axetroy/go-server/core/util"
	"github.com/axetroy/go-server/core/validator"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type SendEmailAuthCodeParams struct {
	Email string `json:"email" valid:"required~请输入邮箱"`
}

type SendPhoneAuthCodeParams struct {
	Phone string `json:"phone" valid:"required~请输入手机号"`
}

func GenerateAuthCode() string {
	return util.RandomString(6)
}

// 发送邮箱验证码 (不需要登陆)
func SendEmailAuthCode(c controller.Context, input SendEmailAuthCodeParams) (res schema.Response) {
	var (
		err error
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

		helper.Response(&res, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	// 生成验证码
	activationCode := GenerateAuthCode()

	// 缓存验证码到 redis
	if err = redis.ClientAuthEmailCode.Set(activationCode, input.Email, time.Minute*10).Err(); err != nil {
		return
	}

	e := email.NewMailer()

	// send email
	if err = e.SendAuthEmail(input.Email, activationCode); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ClientAuthEmailCode.Del(activationCode).Err()
		return
	}

	return
}

// 发送手机验证码 (不需要登陆)
func SendPhoneAuthCode(c controller.Context, input SendPhoneAuthCodeParams) (res schema.Response) {
	var (
		err error
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

		helper.Response(&res, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	// 生成验证码
	activationCode := GenerateAuthCode()

	// 缓存验证码到 redis
	if err = redis.ClientAuthPhoneCode.Set(activationCode, input.Phone, time.Minute*10).Err(); err != nil {
		return
	}

	if err = telephone.GetClient().SendAuthCode(input.Phone, activationCode); err != nil {
		// 如果发送失败，则删除
		_ = redis.ClientAuthPhoneCode.Del(activationCode).Err()
		return
	}

	return
}

func SendEmailAuthCodeRouter(c *gin.Context) {
	var (
		input SendEmailAuthCodeParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = SendEmailAuthCode(controller.NewContext(c), input)
}

func SendPhoneAuthCodeRouter(c *gin.Context) {
	var (
		input SendPhoneAuthCodeParams
		err   error
		res   = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = SendPhoneAuthCode(controller.NewContext(c), input)
}
