// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/middleware"
	"github.com/axetroy/go-server/module/admin"
	"github.com/axetroy/go-server/module/admin/admin_model"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type UpdateProfileParams struct {
	Nickname *string            `json:"nickname" valid:"length(1|36)~昵称长度为1-36位"`
	Gender   *user_model.Gender `json:"gender"`
	Avatar   *string            `json:"avatar"`
}

func GetProfile(context schema.Context) (res schema.Response) {
	var (
		err  error
		data user_schema.Profile
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
				err = common_error.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	user := user_model.User{Id: context.Uid}

	if err = tx.Last(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
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

func GetProfileByAdmin(context schema.Context, userId string) (res schema.Response) {
	var (
		err  error
		data user_schema.Profile
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
				err = common_error.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	tx = database.Db.Begin()

	adminInfo := admin_model.Admin{
		Id: context.Uid,
	}

	if err = tx.Last(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	user := user_model.User{Id: userId}

	if err = tx.Last(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
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

func UpdateProfile(context schema.Context, input UpdateProfileParams) (res schema.Response) {
	var (
		err          error
		data         user_schema.Profile
		tx           *gorm.DB
		shouldUpdate bool
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
				err = common_error.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = common_error.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	updated := user_model.User{}

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
		if err = tx.Table(updated.TableName()).Where(user_model.User{Id: context.Uid}).Updates(updated).Error; err != nil {
			return
		}
	}

	userInfo := user_model.User{
		Id: context.Uid,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
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

func UpdateProfileByAdmin(context schema.Context, userId string, input UpdateProfileParams) (res schema.Response) {
	var (
		err          error
		data         user_schema.Profile
		tx           *gorm.DB
		shouldUpdate bool
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
				err = common_error.ErrUnknown
			}
		}

		if tx != nil {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}

		if err != nil {
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}
	}()

	// 参数校验
	if isValidInput, err = govalidator.ValidateStruct(input); err != nil {
		return
	} else if isValidInput == false {
		err = common_error.ErrInvalidParams
		return
	}

	tx = database.Db.Begin()

	// 检查是不是管理员
	adminInfo := admin_model.Admin{
		Id: context.Uid,
	}

	if err = tx.First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = admin.ErrAdminNotExist
		}
		return
	}

	updated := user_model.User{}

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
		if err = tx.Table(updated.TableName()).Where(user_model.User{Id: userId}).Updates(updated).Error; err != nil {
			return
		}
	}

	userInfo := user_model.User{
		Id: userId,
	}

	if err = tx.First(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = common_error.ErrUserNotExist
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

func GetProfileRouter(ctx *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	res = GetProfile(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	})
}

func GetProfileByAdminRouter(ctx *gin.Context) {
	var (
		err error
		res = schema.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		ctx.JSON(http.StatusOK, res)
	}()

	userId := ctx.Param("user_id")

	res = GetProfileByAdmin(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, userId)
}

func UpdateProfileRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	res = UpdateProfile(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, input)
}

func UpdateProfileByAdminRouter(ctx *gin.Context) {
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
		ctx.JSON(http.StatusOK, res)
	}()

	userId := ctx.Param("user_id")

	if err = ctx.ShouldBindJSON(&input); err != nil {
		err = common_error.ErrInvalidParams
		return
	}

	res = UpdateProfileByAdmin(schema.Context{
		Uid: ctx.GetString(middleware.ContextUidField),
	}, userId, input)
}
