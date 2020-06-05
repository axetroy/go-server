// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package admin

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
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"time"
)

type SignInParams struct {
	Username string `json:"username" validate:"required,min=1,max=36" comment:"帐号"`
	Password string `json:"password" validate:"required,min=6,max=36" comment:"密码"`
	Duration *int64 `json:"duration" validate:"omitempty,number,gt=0" comment:"有效时间"`
}

func Login(input SignInParams) (res schema.Response) {
	var (
		err  error
		data = schema.AdminProfileWithToken{}
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

	tx = database.Db.Begin()

	adminInfo := model.Admin{
		Username: input.Username,
		Password: util.GeneratePassword(input.Password),
	}

	if err = database.Db.Where(&adminInfo).First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
		return
	}

	if err = mapstructure.Decode(adminInfo, &data.AdminProfilePure); err != nil {
		return
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	var duration time.Duration

	if input.Duration != nil {
		duration = time.Duration(*input.Duration * int64(time.Second))
	}

	// generate token
	if t, er := token.Generate(adminInfo.Id, token.StateAdmin, duration); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	return
}

var LoginRouter = router.Handler(func(c router.Context) {
	var (
		input SignInParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return Login(input)
	})
})
