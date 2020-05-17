// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type UpdateProfileParams struct {
	Username *string                    `json:"username"` // 用户名，部分用户有机会修改自己的用户名，比如微信注册的帐号
	Nickname *string                    `json:"nickname" valid:"length(1|36)~昵称长度为1-36位"`
	Gender   *model.Gender              `json:"gender"`
	Avatar   *string                    `json:"avatar"`
	Wechat   *UpdateWechatProfileParams `json:"wechat"` // 更新微信绑定的帐号相关
}

// 绑定的微信信息帐号相关
type UpdateWechatProfileParams struct {
	Nickname  *string `json:"nickname"`   // 用户昵称
	AvatarUrl *string `json:"avatar_url"` // 用户头像
	Gender    *int    `json:"gender"`     // 性别
	Country   *string `json:"country"`    // 国家
	Province  *string `json:"province"`   // 省份
	City      *string `json:"city"`       // 城市
	Language  *string `json:"language"`   // 语言
}

func GetProfile(c helper.Context) (res schema.Response) {
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

		helper.Response(&res, data, nil, err)
	}()

	tx = database.Db.Begin()

	userInfo := model.User{Id: c.Uid}

	if err = tx.Where(&userInfo).Preload("Wechat").Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
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

	return
}

func UpdateProfile(c helper.Context, input UpdateProfileParams) (res schema.Response) {
	var (
		err          error
		data         schema.Profile
		tx           *gorm.DB
		shouldUpdate bool
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

	tx = database.Db.Begin()

	updated := model.User{}

	if input.Username != nil {
		shouldUpdate = true

		if err = validator.ValidateUsername(*input.Username); err != nil {
			return
		}

		u := model.User{Id: c.Uid}

		if err = tx.Where(&u).First(&u).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.UserNotExist
			}
			return
		}

		// 如果没有剩余的重命名次数的话
		if u.UsernameRenameRemaining <= 0 {
			err = exception.RenameUserNameFail
			return
		}

		updated.Username = *input.Username
		updated.UsernameRenameRemaining = u.UsernameRenameRemaining - 1
	}

	if input.Nickname != nil {
		updated.Nickname = input.Nickname
		shouldUpdate = true
	}

	if input.Avatar != nil {
		updated.Avatar = *input.Avatar
		shouldUpdate = true
	}

	if input.Gender != nil {
		updated.Gender = *input.Gender
		shouldUpdate = true
	}

	if shouldUpdate {
		if err = tx.Table(updated.TableName()).Where(model.User{Id: c.Uid}).Updates(updated).Error; err != nil {
			return
		}
	}

	userInfo := model.User{
		Id: c.Uid,
	}

	if err = tx.Where(&userInfo).First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if input.Wechat != nil {
		wechatInfo := model.WechatOpenID{
			Uid: userInfo.Id,
		}
		// 判断该用户是否绑定了微信帐号
		if err = tx.Where(&wechatInfo).First(&wechatInfo).Error; err != nil {
			// 如果没有找到，说明帐号没有绑定微信，抛出异常
			if err == gorm.ErrRecordNotFound {
				err = exception.InvalidParams
			}
			return
		}

		// 更新对应的字段
		wechatUpdated := model.WechatOpenID{}
		shouldUpdateWechat := false

		if input.Wechat.Nickname != nil {
			wechatUpdated.Nickname = input.Wechat.Nickname
			shouldUpdateWechat = true
		}

		if input.Wechat.AvatarUrl != nil {
			wechatUpdated.AvatarUrl = input.Wechat.AvatarUrl
			shouldUpdateWechat = true
		}

		if input.Wechat.Gender != nil {
			wechatUpdated.Gender = input.Wechat.Gender
			shouldUpdateWechat = true
		}

		if input.Wechat.Country != nil {
			wechatUpdated.Country = input.Wechat.Country
			shouldUpdateWechat = true
		}

		if input.Wechat.Province != nil {
			wechatUpdated.Province = input.Wechat.Province
			shouldUpdateWechat = true
		}

		if input.Wechat.City != nil {
			wechatUpdated.City = input.Wechat.City
			shouldUpdateWechat = true
		}

		if input.Wechat.Language != nil {
			wechatUpdated.Language = input.Wechat.Language
			shouldUpdateWechat = true
		}

		if shouldUpdateWechat {
			info := model.WechatOpenID{Id: wechatInfo.Id}
			if err = tx.Where(&info).Updates(wechatUpdated).Error; err != nil {
				if err == gorm.ErrRecordNotFound {
					err = exception.InvalidParams
				}
				return
			}

			wechat := schema.WechatBindingInfo{}

			if err = mapstructure.Decode(info, &wechat); err != nil {
				return
			}

			data.Wechat = &wechat
		}
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var GetProfileRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return GetProfile(helper.NewContext(&c))
	})
})

var UpdateProfileRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateProfileParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return UpdateProfile(helper.NewContext(&c), input)
	})
})
