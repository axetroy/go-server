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

func GetProfileByAdmin(c helper.Context, userId string) (res schema.Response) {
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

	adminInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.Last(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	user := model.User{Id: userId}

	if err = tx.Last(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if err = mapstructure.Decode(user, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = user.PayPassword != nil && len(*user.PayPassword) != 0
	data.CreatedAt = user.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = user.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func UpdateProfileByAdmin(c helper.Context, userId string, input UpdateProfileParams) (res schema.Response) {
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

	// 检查是不是管理员
	adminInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	updated := model.User{}

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
		if err = tx.Table(updated.TableName()).Where(model.User{Id: userId}).Updates(updated).Error; err != nil {
			return
		}
	}

	userInfo := model.User{
		Id: userId,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var GetProfileByAdminRouter = router.Handler(func(c router.Context) {
	userId := c.Param("user_id")

	c.ResponseFunc(nil, func() schema.Response {
		return GetProfileByAdmin(helper.NewContext(&c), userId)
	})
})

var UpdateProfileByAdminRouter = router.Handler(func(c router.Context) {
	var (
		input UpdateProfileParams
	)

	userId := c.Param("user_id")

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return UpdateProfileByAdmin(helper.NewContext(&c), userId, input)
	})
})
