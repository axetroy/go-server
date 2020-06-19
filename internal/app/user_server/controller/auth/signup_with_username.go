// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package auth

import (
	"errors"
	"github.com/axetroy/go-server/internal/app/user_server/controller/wallet"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"time"
)

type SignUpWithUsernameParams struct {
	Username   string  `json:"username" validate:"required,max=32" comment:"用户名"`    // 用户名
	Password   string  `json:"password" validate:"required,max=32" comment:"密码"`     // 密码
	InviteCode *string `json:"invite_code" validate:"omitempty,len=8" comment:"邀请码"` // 邀请码
}

// 创建用户帐号，包括创建的邀请码，钱包数据等，继承到一起
func CreateUserTx(tx *gorm.DB, userInfo *model.User, inviterCode *string) (err error) {
	var (
		newTx bool
	)
	if tx == nil {
		tx = database.Db.Begin()
		newTx = true
	}

	defer func() {
		if newTx {
			if err != nil {
				_ = tx.Rollback().Error
			} else {
				err = tx.Commit().Error
			}
		}
	}()

	if err = tx.Create(userInfo).Error; err != nil {
		return err
	}

	if inviterCode != nil && len(*inviterCode) > 0 {

		inviter := model.User{
			InviteCode: *inviterCode,
		}

		if err := tx.Where(&inviter).Find(&inviter).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = exception.InvalidInviteCode
			}
			return err
		}

		// 如果存在邀请者的话，写入邀请列表中
		if inviter.Id != "" {
			inviteHistory := model.InviteHistory{
				Inviter:       inviter.Id,
				Invitee:       userInfo.Id,
				Status:        model.StatusInviteRegistered,
				RewardSettled: false,
			}

			// 创建邀请记录
			if err = tx.Create(&inviteHistory).Error; err != nil {
				return err
			}
		}
	}

	// 创建用户对应的钱包账号
	for _, walletName := range model.Wallets {
		if err = tx.Table(wallet.GetTableName(walletName)).Create(&model.Wallet{
			Id:       userInfo.Id,
			Currency: walletName,
			Balance:  0,
			Frozen:   0,
		}).Error; err != nil {
			return err
		}
	}

	return nil
}

// 使用用户名注册
func SignUpWithUsername(input SignUpWithUsernameParams) (res schema.Response) {
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

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	if err = validator.ValidateUsername(input.Username); err != nil {
		return
	}

	tx = database.Db.Begin()

	u := model.User{Username: input.Username}

	if err = tx.Model(u).Where("username = ?", input.Username).Find(&u).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
	}

	if u.Id != "" {
		err = exception.UserExist
		return
	}

	userInfo := model.User{
		Username: input.Username,
		Nickname: &input.Username,
		Password: util.GeneratePassword(input.Password),
		Status:   model.UserStatusInit,
		Role:     pq.StringArray{model.DefaultUser.Name},
		Phone:    nil,
		Email:    nil,
		Gender:   model.GenderUnknown,
	}

	if err = CreateUserTx(tx, &userInfo, input.InviteCode); err != nil {
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

var SignUpWithUsernameRouter = router.Handler(func(c router.Context) {
	var (
		input SignUpWithUsernameParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignUpWithUsername(input)
	})
})
