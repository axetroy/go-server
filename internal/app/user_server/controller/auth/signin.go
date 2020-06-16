// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package auth

import (
	"errors"
	"github.com/axetroy/go-server/internal/app/user_server/controller/wallet"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/authentication"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/dotenv"
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/service/wechat"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"
)

type SignInParams struct {
	Account  string `json:"account" validate:"required,max=36" comment:"帐号"`
	Password string `json:"password" validate:"required,max=32" comment:"密码"`
	Duration *int64 `json:"duration" validate:"omitempty,number,gt=0" comment:"有效时间"`
}

type SignInWithEmailParams struct {
	Email    string `json:"email" validate:"required,email" comment:"邮箱"`
	Code     string `json:"code" validate:"required" comment:"验证码"`
	Duration *int64 `json:"duration" validate:"omitempty,number,gt=0" comment:"有效时间"`
}

type SignInWithPhoneParams struct {
	Phone    string `json:"phone" validate:"required,numeric,len=11" comment:"手机号"`
	Code     string `json:"code" validate:"required" comment:"验证码"`
	Duration *int64 `json:"duration" validate:"omitempty,number,gt=0" comment:"有效时间"`
}

type SignInWithWechatParams struct {
	Code     string `json:"code" validate:"required" comment:"微信授权代码"` // 微信小程序授权之后返回的 code
	Duration *int64 `json:"duration" validate:"omitempty,number,gt=0" comment:"有效时间"`
}

type SignInWithOAuthParams struct {
	Code     string `json:"code" validate:"required" comment:"oAuth授权代码"` // oAuth 授权之后回调返回的 code
	Duration *int64 `json:"duration" validate:"omitempty,number,gt=0" comment:"有效时间"`
}

type WechatCompleteParams struct {
	Code     string  `json:"code" validate:"required" comment:"微信授权代码"`               // 微信小程序授权之后返回的 code
	Phone    *string `json:"phone" validate:"omitempty,numeric,len=11" comment:"手机号"` // 手机号
	Username *string `json:"username" validate:"omitempty,max=32" comment:"用户名"`      // 用户名
}

// 普通帐号登陆
func SignIn(c helper.Context, input SignInParams) (res schema.Response) {
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

	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	userInfo := model.User{
		Password: util.GeneratePassword(input.Password),
	}

	if validator.IsPhone(input.Account) {
		// 用手机号登陆
		userInfo.Phone = &input.Account
	} else if validator.IsEmail(input.Account) {
		// 用邮箱登陆
		userInfo.Email = &input.Account
	} else {
		// 用用户名
		userInfo.Username = input.Account
	}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Preload("Wechat").Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
		return
	}

	// 检查用户登录状态
	go func() {
		if er := message_queue.PublishCheckUserLogin(userInfo.Id); er != nil {
			log.Println("检查用户状态失败", c.Uid)
		}
	}()

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

	email, err := redis.ClientAuthEmailCode.Get(input.Code).Result()

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

// 手机 + 验证码登陆
func SignInWithPhone(c helper.Context, input SignInWithPhoneParams) (res schema.Response) {
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

	phone, err := redis.ClientAuthPhoneCode.Get(input.Code).Result()

	// 校验验证码是否正确
	if err != nil || phone != input.Phone {
		err = exception.InvalidParams
	}

	userInfo := model.User{
		Email: &input.Phone,
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

// 使用微信小程序登陆
func SignInWithWechat(c helper.Context, input SignInWithWechatParams) (res schema.Response) {
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

	wechatInfo, wechatErr := wechat.FetchOpenID(input.Code)

	if wechatErr != nil {
		err = wechatErr
		return
	}

	if len(wechatInfo.OpenID) == 0 {
		// 为了测试方便，那么我们就在测试环境下赋值一个假 ID， 否则会报错
		if dotenv.Test {
			wechatInfo.OpenID = "test open id"
		} else {
			err = exception.NoData
			return
		}

	}

	tx = database.Db.Begin()

	wechatOpenID := model.WechatOpenID{
		Id: wechatInfo.OpenID,
	}

	// 去查表
	result := tx.Where(&wechatOpenID).Preload("User").First(&wechatOpenID)

	var userInfo *model.User

	if result.RecordNotFound() {
		var (
			uid      = util.GenerateId()
			username = "v" + uid
		)

		userInfo = &model.User{
			Username:                username,
			Nickname:                &username,
			Password:                util.GeneratePassword(uid),
			Status:                  model.UserStatusInit,
			Role:                    pq.StringArray{model.DefaultUser.Name},
			Gender:                  model.GenderUnknown,
			WechatOpenID:            &wechatOpenID.Id,
			UsernameRenameRemaining: 1, // 允许微信注册的用户可以重命名一次
		}

		if err = tx.Create(userInfo).Error; err != nil {
			return
		}

		if err = tx.Create(&model.WechatOpenID{
			Id:  wechatInfo.OpenID,
			Uid: userInfo.Id,
		}).Error; err != nil {
			return
		}

		// 创建用户对应的钱包账号
		for _, walletName := range model.Wallets {
			if err = tx.Table(wallet.GetTableName(walletName)).Create(&model.Wallet{
				Id:       userInfo.Id,
				Currency: walletName,
				Balance:  0,
				Frozen:   0,
			}).Error; err != nil {
				return
			}
		}

	} else {
		userInfo = &wechatOpenID.User
	}

	if userInfo == nil {
		err = exception.NoData
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	wechatBindingInfo := schema.WechatBindingInfo{}

	if err = mapstructure.Decode(wechatOpenID, &wechatBindingInfo); err != nil {
		return
	}

	data.Wechat = &wechatBindingInfo
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
		Type:    model.LoginLogTypeWechat,          // 微信登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  c.UserAgent,                       // 用户的 userAgent
		LastIp:  c.Ip,                              // 用户的IP
	}

	if err = tx.Create(&loginLog).Error; err != nil {
		return
	}

	return
}

// 使用 oAuth 认证方式登陆
func SignInWithOAuth(c helper.Context, input SignInWithOAuthParams) (res schema.Response) {
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

	uid, err := redis.ClientOAuthCode.Get(input.Code).Result()

	if err != nil {
		return
	}

	var userInfo = model.User{
		Id: uid,
	}

	if err = tx.Where(&userInfo).Preload("Wechat").Find(&userInfo).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	if userInfo.WechatOpenID != nil {
		wechatBindingInfo := schema.WechatBindingInfo{}

		if err = mapstructure.Decode(userInfo.Wechat, &wechatBindingInfo); err != nil {
			return
		}

		data.Wechat = &wechatBindingInfo
	}

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

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var SignInRouter = router.Handler(func(c router.Context) {
	var (
		input SignInParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignIn(helper.NewContext(&c), input)
	})
})

var SignInWithEmailRouter = router.Handler(func(c router.Context) {
	var (
		input SignInWithEmailParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignInWithEmail(helper.NewContext(&c), input)
	})
})

var SignInWithPhoneRouter = router.Handler(func(c router.Context) {
	var (
		input SignInWithPhoneParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignInWithPhone(helper.NewContext(&c), input)
	})
})

var SignInWithWechatRouter = router.Handler(func(c router.Context) {
	var (
		input SignInWithWechatParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignInWithWechat(helper.NewContext(&c), input)
	})
})

var SignInWithOAuthRouter = router.Handler(func(c router.Context) {
	var (
		input SignInWithOAuthParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignInWithOAuth(helper.NewContext(&c), input)
	})
})
