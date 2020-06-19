// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package phone

import (
	"context"
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/service/telephone"
	"time"
)

type SendPhoneCodeForLoginParams struct {
	Phone string `json:"phone" validate:"required,numeric,len=11" comment:"手机号"`
}

// 发送手机验证码 (不需要登陆)
func SendPhoneCodeForLogin(c helper.Context, input SendPhoneCodeForLoginParams) (res schema.Response) {
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
	if err = redis.ClientAuthPhoneCode.Set(context.Background(), activationCode, input.Phone, time.Minute*10).Err(); err != nil {
		return
	}

	if err = telephone.GetClient().SendAuthCode(input.Phone, activationCode); err != nil {
		// 如果发送失败，则删除
		_ = redis.ClientAuthPhoneCode.Del(context.Background(), activationCode).Err()
		return
	}

	return
}

var SendPhoneCodeForLoginRouter = router.Handler(func(c router.Context) {
	var (
		input SendPhoneCodeForLoginParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SendPhoneCodeForLogin(helper.NewContext(&c), input)
	})
})
