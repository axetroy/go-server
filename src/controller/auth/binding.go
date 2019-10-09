package auth

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/service/redis"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type BindingEmailParams struct {
	Email string `json:"email" valid:"required~请输入邮箱"` // 邮箱
	Code  string `json:"code" valid:"required~请输入验证码"` // 邮箱收到的验证码
}

type BindingPhoneParams struct {
	Phone string `json:"phone" valid:"required~请输入手机号"` // 手机号
	Code  string `json:"code" valid:"required~请输入验证码"`  // 手机收到的验证码
}

type BindingWechatMiniAppParams struct {
	Code string `json:"code" valid:"required~请输入微信认证码"` // 微信小程序调用 wx.login() 之后，返回的 code
}

// 绑定邮箱
func BindingEmail(c controller.Context, input BindingEmailParams) (res schema.Response) {
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
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
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
		err = exception.DuplicateBinding
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
		err = exception.DuplicateBinding
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
func BindingPhone(c controller.Context, input BindingPhoneParams) (res schema.Response) {
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
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
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
		err = exception.DuplicateBinding
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
		err = exception.DuplicateBinding
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
func BindingWechat(c controller.Context, input BindingWechatMiniAppParams) (res schema.Response) {
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
		return
	} else if isValidInput == false {
		err = exception.InvalidParams
		return
	}

	wechatRes, err := FetchWechatInfo(input.Code)

	if err != nil {
		return
	}

	tx = database.Db.Begin()

	wechatInfo := model.WechatOpenID{
		Id: wechatRes.OpenID,
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
		err = exception.DuplicateBinding
		return
	}

	// 乳沟没有绑定，则更新绑定信息
	if err = tx.Where(&wechatInfo).Update("uid", c.Context).Error; err != nil {
		return
	}

	return
}

func BindingEmailRouter(c *gin.Context) {
	var (
		input BindingEmailParams
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

	res = BindingEmail(controller.NewContext(c), input)
}

func BindingPhoneRouter(c *gin.Context) {
	var (
		input BindingPhoneParams
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

	res = BindingPhone(controller.NewContext(c), input)
}

func BindingWechatRouter(c *gin.Context) {
	var (
		input BindingWechatMiniAppParams
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

	res = BindingWechat(controller.NewContext(c), input)
}
