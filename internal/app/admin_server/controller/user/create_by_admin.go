// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package user

import (
	"encoding/json"
	"errors"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/wallet"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/message_queue"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"time"
)

type CreateUserParams struct {
	Username *string      `json:"username"`
	Email    *string      `json:"email"`
	Phone    *string      `json:"phone"`
	Nickname *string      `json:"nickname"`
	Gender   model.Gender `json:"gender"`
	Password string       `json:"password"`
}

func CreateUser(input CreateUserParams) (res schema.Response) {
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

	if input.Password == "" {
		err = exception.RequirePassword
		return
	}

	if input.Username == nil && input.Phone == nil && input.Email == nil {
		err = errors.New("请输入账号")
		return
	}

	tx = database.Db.Begin()

	var (
		username string
		uid      = util.GenerateId()
	)

	if input.Username == nil {
		username = "u" + uid
	} else {
		username = *input.Username
	}

	var (
		existUserInfo = model.User{}
	)

	if input.Username != nil {
		if err = tx.Where("username = ?", *input.Username).Find(&existUserInfo).Error; err != nil {
			// 如果找不到这个用户
			// 说明用户没存在
			if err != gorm.ErrRecordNotFound {
				return
			}
		}

		if existUserInfo.Id != "" {
			err = exception.UserExist
			return
		}
	}

	if input.Email != nil {
		if err = tx.Where("email = ?", *input.Email).Find(&existUserInfo).Error; err != nil {
			// 如果找不到这个用户
			// 说明用户没存在
			if err != gorm.ErrRecordNotFound {
				return
			}
		}

		if existUserInfo.Id != "" {
			err = exception.UserExist
			return
		}
	}

	if input.Phone != nil {
		if err = tx.Where("phone = ?", *input.Phone).Find(&existUserInfo).Error; err != nil {
			// 如果找不到这个用户
			// 说明用户没存在
			if err != gorm.ErrRecordNotFound {
				return
			}
		}

		if existUserInfo.Id != "" {
			err = exception.UserExist
			return
		}
	}

	nickname := &username

	if input.Nickname != nil {
		nickname = input.Nickname
	}

	userInfo := model.User{
		Username: username,
		Nickname: nickname,
		Password: util.GeneratePassword(input.Password),
		Status:   model.UserStatusInactivated, // 开始时未激活状态
		Role:     pq.StringArray{model.DefaultUser.Name},
		Phone:    input.Phone,
		Email:    input.Email,
		Gender:   input.Gender,
	}

	if err = tx.Create(&userInfo).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// 创建用户对应的钱包账号
	for _, walletName := range model.Wallets {
		if err = tx.Table(wallet.GetTableName(walletName)).Create(&model.Wallet{
			Id:       userInfo.Id,
			Currency: walletName,
			Balance:  0,
			Frozen:   0,
		}).Error; err != nil {
			return
		}
	}

	// 如果是以邮箱注册的，那么发送激活链接
	if userInfo.Email != nil && len(*userInfo.Email) != 0 {
		// 生成激活码
		activationCode := "activation-" + userInfo.Id

		// 把激活码存到 redis
		if err = redis.ClientActivationCode.Set(activationCode, userInfo.Id, time.Minute*30).Err(); err != nil {
			return
		}

		// 把 "发送激活码" 加入消息队列
		var body []byte

		if body, err = json.Marshal(message_queue.BodySendActivationEmail{
			Email: *input.Email,
			Code:  activationCode,
		}); err != nil {
			return
		}

		if err = message_queue.Publish(message_queue.TopicSendEmail, body); err != nil {
			return
		}

		return
	}
	return
}

var CreateUserRouter = router.Handler(func(c router.Context) {
	var (
		input CreateUserParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return CreateUser(input)
	})
})
