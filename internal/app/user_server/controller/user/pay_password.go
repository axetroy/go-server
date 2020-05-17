// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/email"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/service/telephone"
	"github.com/jinzhu/gorm"
	"time"
)

type SetPayPasswordParams struct {
	Password        string `json:"password" valid:"required~请输入密码,int~请输入纯数字的密码,length(6|6)~密码长度为6位"`
	PasswordConfirm string `json:"password_confirm" valid:"required~请输入确认密码,int~请输入纯数字的确认密码,length(6|6)~确认密码长度为6位"`
}

type UpdatePayPasswordParams struct {
	OldPassword string `json:"old_password" valid:"required~请输入旧密码,int~请输入纯数字的旧密码,length(6|6)~旧密码长度为6位"`
	NewPassword string `json:"new_password" valid:"required~请输入新密码,int~请输入纯数字的新密码,length(6|6)~新密码长度为6位"`
}

type ResetPayPasswordParams struct {
	Code        string `json:"code" valid:"required~请输入重置码"`                                                // 重置码
	NewPassword string `json:"new_password" valid:"required~请输入新的交易密码,int~请输入纯数字的旧密码,length(6|6)~新密码长度为6位"` // 新的交易密码
}

func GenerateResetPayPasswordCode(uid string) string {
	codeId := "reset-pay-" + util.GenerateId() + uid
	return util.MD5(codeId)
}

func SetPayPassword(c helper.Context, input SetPayPasswordParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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

		helper.Response(&res, nil, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if input.Password != input.PasswordConfirm {
		err = exception.InvalidConfirmPassword
		return
	}

	userInfo := model.User{Id: c.Uid}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if userInfo.PayPassword != nil {
		err = exception.PayPasswordSet
		return
	}

	newPassword := util.GeneratePassword(input.Password)

	// 更新交易密码
	if err = database.Db.Model(userInfo).Update("pay_password", newPassword).Error; err != nil {
		return
	}

	return
}

var SetPayPasswordRouter = router.Handler(func(c router.Context) {
	var (
		input SetPayPasswordParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SetPayPassword(helper.NewContext(&c), input)
	})
})

func UpdatePayPassword(c helper.Context, input UpdatePayPasswordParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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

		helper.Response(&res, nil, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if input.OldPassword == input.NewPassword {
		err = exception.PasswordDuplicate
		return
	}

	userInfo := model.User{Id: c.Uid}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 如果还没有设置过交易密码，就不会有更新
	if userInfo.PayPassword == nil {
		err = exception.RequirePayPasswordSet
		return
	}

	oldPwd := util.GeneratePassword(input.OldPassword)

	// 旧密码不匹配
	if *userInfo.PayPassword != oldPwd {
		err = exception.InvalidPassword
		return
	}

	newPwd := util.GeneratePassword(input.NewPassword)

	// 更新交易密码
	if err = database.Db.Model(userInfo).Update("pay_password", newPwd).Error; err != nil {
		return
	}

	return
}

var UpdatePayPasswordRouter = router.Handler(func(c router.Context) {
	var (
		input UpdatePayPasswordParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return UpdatePayPassword(helper.NewContext(&c), input)
	})
})

func SendResetPayPassword(c helper.Context) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
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

		helper.Response(&res, nil, nil, err)
	}()

	userInfo := model.User{Id: c.Uid}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 生成重置码
	var resetCode = GenerateResetPayPasswordCode(userInfo.Id)

	// redis缓存重置码
	if err = redis.ClientResetCode.Set(resetCode, userInfo.Id, time.Minute*10).Err(); err != nil {
		return
	}

	if userInfo.Email != nil {
		// 发送邮件
		go func() {
			e := email.NewMailer()
			_ = e.SendForgotTradePasswordEmail(*userInfo.Email, resetCode)
		}()
	} else if userInfo.Phone != nil {
		go func() {
			if err = telephone.GetClient().SendResetPasswordCode(*userInfo.Phone, resetCode); err != nil {
				// 如果发送失败，则删除
				_ = redis.ClientAuthPhoneCode.Del(resetCode).Err()
				return
			}
		}()
	} else {
		// 无效的用户
		err = exception.NoData
	}

	return
}

var SendResetPayPasswordRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return SendResetPayPassword(helper.NewContext(&c))
	})
})

func ResetPayPassword(c helper.Context, input ResetPayPasswordParams) (res schema.Response) {
	var (
		err error
		tx  *gorm.DB
		uid string
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

		helper.Response(&res, nil, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	userInfo := model.User{Id: c.Uid}

	tx = database.Db.Begin()

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	// 如果还没有设置过交易密码，就不会有重置
	if userInfo.PayPassword == nil {
		err = exception.RequirePayPasswordSet
		return
	}

	if uid, err = redis.ClientResetCode.Get(input.Code).Result(); err != nil {
		err = exception.InvalidResetCode
		return
	}

	// 即使有了重置码，不是自己的账号也不能用
	if userInfo.Id != uid {
		err = exception.NoPermission
		return
	}

	// 更新交易密码
	if err = database.Db.Model(userInfo).Update("pay_password", input.NewPassword).Error; err != nil {
		return
	}

	// 重置密码之后，删除重置码
	if _, err = redis.ClientResetCode.Del(input.Code).Result(); err != nil {
		return
	}

	return
}

var ResetPayPasswordRouter = router.Handler(func(c router.Context) {
	var (
		input ResetPayPasswordParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return ResetPayPassword(helper.NewContext(&c), input)
	})
})
