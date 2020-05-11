package auth

import (
	"errors"
	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/helper"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/service/wechat"
	"github.com/axetroy/go-server/internal/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type UnbindingEmailParams struct {
	Code string `json:"code" valid:"required~请输入验证码"` // 解除邮箱绑定前，需要发送邮箱验证码验证
}

type UnbindingPhoneParams struct {
	Code string `json:"code" valid:"required~请输入验证码"` // 解除手机号绑定前，需要发送手机验证码验证
}

type UnbindingWechatParams struct {
	Code string `json:"code" valid:"required~请输入验证码"` // 验证码，如果帐号已绑定手机，则为手机号收到的验证码，如果有为邮箱，则用邮箱收到的验证码，否则使用 `wx.login()` 返回的 code
}

// 解除邮箱绑定
func UnbindingEmail(c controller.Context, input UnbindingEmailParams) (res schema.Response) {
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

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	userInfo := model.User{}

	if err = tx.Where(&userInfo).Find(&userInfo).Error; err != nil {
		return
	}

	// 如果邮箱为空，则不需要解绑
	if userInfo.Email == nil {
		err = exception.NoData
		return
	}

	// 校验验证码正确不正确
	email, err := redis.ClientAuthEmailCode.Get(input.Code).Result()

	if err != nil {
		return
	}

	// 如果邮箱不匹配，则校验失败
	if email != *userInfo.Email {
		err = exception.InvalidParams
		return
	}

	if err = tx.Where(model.User{Id: c.Uid}).Update("email", nil).Error; err != nil {
		return
	}

	return
}

// 解除手机绑定
func UnbindingPhone(c controller.Context, input UnbindingPhoneParams) (res schema.Response) {
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

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	userInfo := model.User{}

	if err = tx.Where(&userInfo).Find(&userInfo).Error; err != nil {
		return
	}

	// 如果邮箱为空，则不需要解绑
	if userInfo.Phone == nil {
		err = exception.NoData
		return
	}

	// 校验验证码正确不正确
	phone, err := redis.ClientAuthPhoneCode.Get(input.Code).Result()

	if err != nil {
		return
	}

	// 如果邮箱不匹配，则校验失败
	if phone != *userInfo.Phone {
		err = exception.InvalidParams
		return
	}

	if err = tx.Where(model.User{Id: c.Uid}).Update("phone", nil).Error; err != nil {
		return
	}

	return
}

// 解除微信绑定
func UnbindingWechat(c controller.Context, input UnbindingWechatParams) (res schema.Response) {
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

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	tx = database.Db.Begin()

	userInfo := model.User{}

	if err = tx.Where(&userInfo).Find(&userInfo).Error; err != nil {
		return
	}

	// 校验验证码是否正确
	if userInfo.Phone != nil {
		// 如果用户已有手机号，则用手机号作为验证码

		// 校验验证码正确不正确
		phone, err := redis.ClientAuthPhoneCode.Get(input.Code).Result()

		if err != nil {
			return
		}

		// 如果不匹配，则校验失败
		if phone != *userInfo.Phone {
			err = exception.InvalidParams
			return
		}
	} else if userInfo.Email != nil {
		// 	如果用户已有邮箱，则用邮箱作为验证码

		// 校验验证码正确不正确
		email, err := redis.ClientAuthEmailCode.Get(input.Code).Result()

		if err != nil {
			return
		}

		// 如果邮箱不匹配，则校验失败
		if email != *userInfo.Email {
			err = exception.InvalidParams
			return
		}
	} else {
		// 否则按照 `wx.login()` 返回的 code
		weRes, err := wechat.FetchOpenID(input.Code)

		if err != nil {
			return
		}

		wechatInfo := model.WechatOpenID{
			Id: weRes.OpenID,
		}

		if err = tx.Where(&wechatInfo).First(&wechatInfo).Error; err != nil {
			return
		}

		// 如果 UID 对不上，则说明 code 不正确
		if wechatInfo.Uid != c.Uid {
			err = exception.InvalidParams
		}
	}

	wechatInfo := model.WechatOpenID{
		Uid: c.Uid,
	}

	if err = tx.Where(&wechatInfo).First(&wechatInfo).Error; err != nil {
		return
	}

	// 解除绑定
	if err = tx.Model(&wechatInfo).Where("id = ?", wechatInfo.Id).Update("uid", nil).Error; err != nil {
		return
	}

	return
}

func UnbindingEmailRouter(c *gin.Context) {
	var (
		input UnbindingEmailParams
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

	res = UnbindingEmail(controller.NewContext(c), input)
}

func UnbindingPhoneRouter(c *gin.Context) {
	var (
		input UnbindingPhoneParams
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

	res = UnbindingPhone(controller.NewContext(c), input)
}

func UnbindingWechatRouter(c *gin.Context) {
	var (
		input UnbindingWechatParams
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

	res = UnbindingWechat(controller.NewContext(c), input)
}
