package auth

import (
	"context"
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/authentication"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type SignInWithEmailParams struct {
	Email    string `json:"email" validate:"required,email" comment:"邮箱"`
	Code     string `json:"code" validate:"required" comment:"验证码"`
	Duration *int64 `json:"duration" validate:"omitempty,number,gt=0" comment:"有效时间"`
}

// 邮箱 + 验证码登陆
func SignInWithEmail(c helper.Context, input SignInWithEmailParams) (res schema.Response) {
	var (
		err  error
		data = &schema.ProfileWithToken{}
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

		helper.Response(&res, data, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	email, err := redis.ClientAuthEmailCode.Get(context.Background(), input.Code).Result()

	// 校验验证码是否正确
	if err != nil || email != input.Email {
		err = exception.InvalidParams
	}

	userInfo := model.User{
		Email: &input.Email,
	}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Preload("Wechat").Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
		return
	}

	if err = userInfo.CheckStatusValid(); err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	if userInfo.WechatOpenID != nil {
		if err = mapstructure.Decode(userInfo.Wechat, &data.Wechat); err != nil {
			return
		}
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	var duration time.Duration

	if input.Duration != nil {
		duration = time.Duration(*input.Duration * int64(time.Second))
	} else {
		duration = time.Hour * 6
	}

	// generate token
	if t, er := authentication.Gateway(false).Generate(userInfo.Id, duration); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	loginLog := model.LoginLog{
		Uid:     userInfo.Id,                       // 用户ID
		Type:    model.LoginLogTypeUserName,        // 默认用户名登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  c.UserAgent,                       // 用户的 userAgent
		LastIp:  c.Ip,                              // 用户的IP
	}

	if err = tx.Create(&loginLog).Error; err != nil {
		return
	}

	return
}

var SignInWithEmailRouter = router.Handler(func(c router.Context) {
	var (
		input SignInWithEmailParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignInWithEmail(helper.NewContext(&c), input)
	})
})
