// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/helper"
	"github.com/axetroy/go-server/src/middleware"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/validator"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
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

func GetProfile(c controller.Context) (res schema.Response) {
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

func GetProfileByAdmin(c controller.Context, userId string) (res schema.Response) {
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

func UpdateProfile(c controller.Context, input UpdateProfileParams) (res schema.Response) {
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

		helper.Response(&res, data, err)
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

func UpdateProfileByAdmin(c controller.Context, userId string, input UpdateProfileParams) (res schema.Response) {
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

		helper.Response(&res, data, err)
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

func GetProfileRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	res = GetProfile(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	})
}

func GetProfileByAdminRouter(c *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	userId := c.Param("user_id")

	res = GetProfileByAdmin(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, userId)
}

func UpdateProfileRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateProfileParams
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

	res = UpdateProfile(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, input)
}

func UpdateProfileByAdminRouter(c *gin.Context) {
	var (
		err   error
		res   = schema.Response{}
		input UpdateProfileParams
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		c.JSON(http.StatusOK, res)
	}()

	userId := c.Param("user_id")

	if err = c.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = UpdateProfileByAdmin(controller.Context{
		Uid: c.GetString(middleware.ContextUidField),
	}, userId, input)
}
