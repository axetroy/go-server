// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"time"

	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
)

type UpdatePasswordParams struct {
	OldPassword     string `json:"old_password" validate:"required,max=36" comment:"旧密码"`                          // 旧密码
	NewPassword     string `json:"new_password" validate:"required,max=36" comment:"新密码"`                          // 新密码
	ConfirmPassword string `json:"confirm_password" validate:"required,max=36,eqfield=NewPassword" comment:"确认密码"` // 确认密码
}

func UpdatePassword(c helper.Context, input UpdatePasswordParams) (res schema.Response) {
	var (
		err  error
		data schema.AdminProfile
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

	if input.NewPassword != input.ConfirmPassword {
		err = exception.InvalidConfirmPassword
		return
	}

	tx = database.Db.Begin()

	myInfo := model.Admin{
		Id: c.Uid,
	}

	if err = tx.First(&myInfo).Error; err != nil {
		// 没有找到管理员
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
		}
		return
	}

	// 校验密码是否正确
	if myInfo.Password != util.GeneratePassword(input.OldPassword) {
		err = exception.InvalidOldPassword
		return
	}

	newPassword := util.GeneratePassword(input.NewPassword)

	if err = tx.Model(&myInfo).Update("password", newPassword).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
		return
	}

	if err = mapstructure.Decode(myInfo, &data.AdminProfilePure); err != nil {
		return
	}

	if len(data.Accession) == 0 {
		data.Accession = []string{}
	}

	data.CreatedAt = myInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = myInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

var UpdatePasswordRouter = router.Handler(func(c router.Context) {
	var (
		input UpdatePasswordParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return UpdatePassword(helper.NewContext(&c), input)
	})
})
