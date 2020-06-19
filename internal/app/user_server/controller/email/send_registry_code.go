package email

import (
	"context"
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/email"
	"github.com/axetroy/go-server/internal/service/redis"
	"time"
)

type SendEmailCodeForRegistryParams struct {
	Email string `json:"email" validate:"required,email" comment:"邮箱"`
}

// 发送邮箱验证码 (不需要登陆)
func SendEmailCodeForRegistry(c helper.Context, input SendEmailCodeForRegistryParams) (res schema.Response) {
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

		helper.Response(&res, nil, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	// 生成验证码
	activationCode := util.RandomNumeric(6)

	// 缓存验证码到 redis
	if err = redis.ClientAuthEmailCode.Set(context.Background(), activationCode, input.Email, time.Minute*10).Err(); err != nil {
		return
	}

	e, err := email.NewMailer()

	if err != nil {
		return
	}

	// send email
	if err = e.SendAuthEmail(input.Email, activationCode); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ClientAuthEmailCode.Del(context.Background(), activationCode).Err()
		return
	}

	return
}

var SendEmailCodeForRegistryRouter = router.Handler(func(c router.Context) {
	var (
		input SendEmailCodeForRegistryParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SendEmailCodeForRegistry(helper.NewContext(&c), input)
	})
})
