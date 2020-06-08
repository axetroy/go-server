package auth

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/service/wechat"
	"github.com/jinzhu/gorm"
)

type BindingEmailParams struct {
	Email string `json:"email" validate:"required,email,max=36" comment:"邮箱"` // 邮箱
	Code  string `json:"code" validate:"required,len=6" comment:"验证码"`        // 邮箱收到的验证码
}

type BindingPhoneParams struct {
	Phone string `json:"phone" validate:"required,numeric,len=11" comment:"手机号"` // 手机号
	Code  string `json:"code" validate:"required,len=6" comment:"验证码"`           // 手机收到的验证码
}

type BindingWechatMiniAppParams struct {
	Code string `json:"code" validate:"required,max=255" comment:"微信授权代码"` // 微信小程序调用 wx.login() 之后，返回的 code
}

// 绑定邮箱
func BindingEmail(c helper.Context, input BindingEmailParams) (res schema.Response) {
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

	userInfo := model.User{
		Email: &input.Email,
	}

	tx = database.Db.Begin()

	err1 := tx.Where(&userInfo).Last(&userInfo).Error

	if err1 != gorm.ErrRecordNotFound {
		err = err1
		return
	}

	// 如果能找到帐号，说明已经绑定了
	if userInfo.Id != "" {
		err = exception.Duplicate
		return
	}

	userInfo.Id = c.Uid
	userInfo.Email = nil

	err2 := tx.Where(&userInfo).Last(&userInfo).Error

	if err2 != nil {
		err = err2
		return
	}

	// 如果该用户已经绑定过邮箱了
	if userInfo.Email != nil {
		err = exception.Duplicate
		return
	}

	// 校验验证码正确不正确
	email, err := redis.ClientAuthEmailCode.Get(input.Code).Result()

	if err != nil {
		return
	}

	// 验证是否正确
	if email != input.Email {
		err = exception.InvalidParams
		return
	}

	if err = tx.Model(&userInfo).Where("id = ?", c.Uid).Update("email", input.Email).Error; err != nil {
		return
	}

	return
}

// 绑定手机号
func BindingPhone(c helper.Context, input BindingPhoneParams) (res schema.Response) {
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

	userInfo := model.User{
		Phone: &input.Phone,
	}

	tx = database.Db.Begin()

	err1 := tx.Where(&userInfo).Last(&userInfo).Error

	if err1 != gorm.ErrRecordNotFound {
		err = err1
		return
	}

	// 如果能找到帐号，说明已经绑定了
	if userInfo.Id != "" {
		err = exception.Duplicate
		return
	}

	userInfo.Id = c.Uid
	userInfo.Email = nil

	err2 := tx.Where(&userInfo).Last(&userInfo).Error

	if err2 != nil {
		err = err2
		return
	}

	// 如果该用户已经绑定过邮箱了
	if userInfo.Email != nil {
		err = exception.Duplicate
		return
	}

	// 校验验证码正确不正确
	phone, err := redis.ClientAuthPhoneCode.Get(input.Code).Result()

	if err != nil {
		return
	}

	// 验证是否正确
	if phone != input.Phone {
		err = exception.InvalidParams
		return
	}

	if err = tx.Model(&userInfo).Where("id = ?", c.Uid).Update("phone", input.Phone).Error; err != nil {
		return
	}

	return
}

// 绑定微信
func BindingWechat(c helper.Context, input BindingWechatMiniAppParams) (res schema.Response) {
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

	wechatRes, err := wechat.FetchOpenID(input.Code)

	if err != nil {
		return
	}

	tx = database.Db.Begin()

	wechatInfo := model.WechatOpenID{
		Id: wechatRes.OpenID,
	}

	userInfo := model.User{Id: c.Uid}

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		return
	}

	if userInfo.WechatOpenID != nil {
		err = exception.Duplicate
		return
	}

	err = tx.Where(&wechatInfo).Preload("User").First(&wechatInfo).Error

	if err != nil {
		// 如果不存在，说明没有被绑定过，则创建
		if err == gorm.ErrRecordNotFound {
			wechatInfo.Uid = c.Uid

			if err = tx.Create(wechatInfo).Error; err != nil {
				return
			}

			return
		} else {
			return
		}
	}

	// 如果 Uid 不为空，则说明这个微信绑定过帐号了
	if wechatInfo.Uid != "" {
		err = exception.Duplicate
		return
	}

	// 如果没有绑定，则更新绑定信息
	if err = tx.Model(&wechatInfo).Where("id = ?", wechatInfo.Id).Update(model.WechatOpenID{Uid: c.Uid}).Error; err != nil {
		return
	}

	// 更新用户信息
	if err = tx.Model(&userInfo).Where("id = ?", userInfo.Id).Update(model.User{WechatOpenID: &wechatInfo.Id}).Error; err != nil {
		return
	}

	return
}

var BindingEmailRouter = router.Handler(func(c router.Context) {
	var (
		input BindingEmailParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return BindingEmail(helper.NewContext(&c), input)
	})
})

var BindingPhoneRouter = router.Handler(func(c router.Context) {
	var (
		input BindingPhoneParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return BindingPhone(helper.NewContext(&c), input)
	})
})

var BindingWechatRouter = router.Handler(func(c router.Context) {
	var (
		input BindingWechatMiniAppParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return BindingWechat(helper.NewContext(&c), input)
	})
})
