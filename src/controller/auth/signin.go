// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/config"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/wallet"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/service/redis"
	"github.com/axetroy/go-server/src/service/token"
	"github.com/axetroy/go-server/src/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"io/ioutil"
	"net/http"
	"time"
)

type SignInParams struct {
	Account  string `json:"account" valid:"required~请输入登陆账号"`
	Password string `json:"password" valid:"required~请输入密码"`
}

type SignInWithEmailParams struct {
	Email string `json:"email" valid:"required~请输入邮箱"`
	Code  string `json:"code" valid:"required~请输入验证码"`
}

type SignInWithPhoneParams struct {
	Phone string `json:"phone" valid:"required~请输入手机号"`
	Code  string `json:"code" valid:"required~请输入验证码"`
}

type SignInWithWechatParams struct {
	Code string `json:"code" valid:"required~请输入微信授权代码"` // 微信小程序授权之后返回的 code
}

type WechatResponse struct {
	OpenID     string `json:"openid"`      // 用户唯一标识
	SessionKey string `json:"session_key"` // 会话密钥
	UnionID    string `json:"unionid"`     // 用户在开放平台的唯一标识符，在满足 UnionID 下发条件的情况下会返回
	ErrCode    int    `json:"errcode"`     // 错误码
	ErrMsg     string `json:"errmsg"`      // 错误信息
}

type WechatCompleteParams struct {
	Code     string  `json:"code" valid:"required~请输入微信授权代码"` // 微信小程序授权之后返回的 code
	Phone    *string `json:"phone"`                           // 手机号
	Username *string `json:"username"`                        // 用户名
}

// 普通帐号登陆
func SignIn(c controller.Context, input SignInParams) (res schema.Response) {
	var (
		err          error
		data         = &schema.ProfileWithToken{}
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		err = exception.WrapValidatorError(err)
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	userInfo := model.User{
		Password: util.GeneratePassword(input.Password),
	}

	if util.IsPhone(input.Account) {
		// 用手机号登陆
		userInfo.Phone = &input.Account
	} else if govalidator.IsEmail(input.Account) {
		// 用邮箱登陆
		userInfo.Email = &input.Account
	} else {
		// 用用户名
		userInfo.Username = input.Account
	}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
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

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	if t, er := token.Generate(userInfo.Id, false); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	log := model.LoginLog{
		Uid:     userInfo.Id,                       // 用户ID
		Type:    model.LoginLogTypeUserName,        // 默认用户名登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  c.UserAgent,                       // 用户的 userAgent
		LastIp:  c.Ip,                              // 用户的IP
	}

	if err = tx.Create(&log).Error; err != nil {
		return
	}

	return
}

// 邮箱 + 验证码登陆
func SignInWithEmail(c controller.Context, input SignInWithEmailParams) (res schema.Response) {
	var (
		err          error
		data         = &schema.ProfileWithToken{}
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		err = exception.WrapValidatorError(err)
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
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

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
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

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	if t, er := token.Generate(userInfo.Id, false); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	log := model.LoginLog{
		Uid:     userInfo.Id,                       // 用户ID
		Type:    model.LoginLogTypeUserName,        // 默认用户名登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  c.UserAgent,                       // 用户的 userAgent
		LastIp:  c.Ip,                              // 用户的IP
	}

	if err = tx.Create(&log).Error; err != nil {
		return
	}

	return
}

// 手机 + 验证码登陆
func SignInWithPhone(c controller.Context, input SignInWithPhoneParams) (res schema.Response) {
	var (
		err          error
		data         = &schema.ProfileWithToken{}
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		err = exception.WrapValidatorError(err)
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
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

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
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

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	if t, er := token.Generate(userInfo.Id, false); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	log := model.LoginLog{
		Uid:     userInfo.Id,                       // 用户ID
		Type:    model.LoginLogTypeUserName,        // 默认用户名登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  c.UserAgent,                       // 用户的 userAgent
		LastIp:  c.Ip,                              // 用户的IP
	}

	if err = tx.Create(&log).Error; err != nil {
		return
	}

	return
}

func FetchWechatInfo(code string) (*WechatResponse, error) {
	wechatUrl := fmt.Sprintf("https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code", config.Wechat.AppID, config.Wechat.Secret, code)

	r, reqErr := http.Get(wechatUrl)

	if reqErr != nil {
		return nil, reqErr
	}

	resBytes, ioErr := ioutil.ReadAll(r.Body)

	if ioErr != nil {
		return nil, ioErr
	}

	reqRes := WechatResponse{}

	if jsonErr := json.Unmarshal(resBytes, &reqRes); jsonErr != nil {
		return nil, jsonErr
	}

	return &reqRes, nil
}

// 使用微信小程序登陆
func SignInWithWechat(context controller.Context, input SignInWithWechatParams) (res schema.Response) {
	var (
		err          error
		data         = &schema.ProfileWithToken{}
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

		helper.Response(&res, data, err)
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		err = exception.WrapValidatorError(err)
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	wechatInfo, wechatErr := FetchWechatInfo(input.Code)

	if wechatErr != nil {
		err = wechatErr
		return
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
			Username: username,
			Nickname: &username,
			Password: util.GeneratePassword(uid),
			Status:   model.UserStatusInit,
			Role:     pq.StringArray{model.DefaultUser.Name},
			Gender:   model.GenderUnknown,
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

		return
	} else {
		userInfo = &wechatOpenID.User
	}

	if userInfo == nil {
		err = exception.NoData
		return
	}

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	if t, er := token.Generate(userInfo.Id, false); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	log := model.LoginLog{
		Uid:     userInfo.Id,                       // 用户ID
		Type:    model.LoginLogTypeWechat,          // 微信登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  context.UserAgent,                 // 用户的 userAgent
		LastIp:  context.Ip,                        // 用户的IP
	}

	if err = tx.Create(&log).Error; err != nil {
		return
	}

	return
}

func SignInRouter(c *gin.Context) {
	var (
		input SignInParams
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

	res = SignIn(controller.NewContext(c), input)
}

func SignInWithEmailRouter(c *gin.Context) {
	var (
		input SignInWithEmailParams
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

	res = SignInWithEmail(controller.NewContext(c), input)
}

func SignInWithPhoneRouter(c *gin.Context) {
	var (
		input SignInWithPhoneParams
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

	res = SignInWithPhone(controller.NewContext(c), input)
}

func SignInWithWechatRouter(c *gin.Context) {
	var (
		input SignInWithWechatParams
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

	res = SignInWithWechat(controller.NewContext(c), input)
}
