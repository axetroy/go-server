// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth

import (
	"encoding/json"
	"errors"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/module/invite/invite_model"
	"github.com/axetroy/go-server/module/message_queue"
	"github.com/axetroy/go-server/module/role/role_model"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/module/wallet"
	"github.com/axetroy/go-server/module/wallet/wallet_model"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/service/redis"
	"github.com/axetroy/go-server/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type SignUpParams struct {
	Username   *string `json:"username"`
	Email      *string `json:"email"`
	Phone      *string `json:"phone"`
	Password   string  `json:"password"`
	MCode      *string `json:"mcode"`       // 手机验证码
	InviteCode *string `json:"invite_code"` // 邀请码
}

func SignUp(input SignUpParams) (res schema.Response) {
	var (
		err     error
		data    user_schema.Profile
		tx      *gorm.DB
		inviter *user_model.User // 邀请人信息
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = schema.StatusSuccess
		}

	}()

	if input.Password == "" {
		err = common_error.ErrRequirePassword
		return
	}

	if input.Username == nil && input.Phone == nil && input.Email == nil {
		err = errors.New("请输入账号")
		return
	}

	if input.Phone != nil {
		if input.MCode == nil {
			err = errors.New("请输入短信验证码")
			return
		}

		// TODO: 校验手机号码是否正确

		// TODO: 验证短信验证码是否正确
		// 因为现在没有引入短信服务, 所以暂时没有这一块的功能
	}

	tx = database.Db.Begin()

	var (
		username string
		uid      = util.GenerateId()
	)

	if input.Username == nil {
		username = "用户" + uid
	} else {
		username = *input.Username
	}

	var (
		existUserInfo = user_model.User{}
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
			err = common_error.ErrUserExist
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
			err = common_error.ErrUserExist
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
			err = common_error.ErrUserExist
			return
		}
	}

	// 填入了邀请码，则去校验邀请码是否正确
	if input.InviteCode != nil && *input.InviteCode != "" {
		u := user_model.User{
			InviteCode: *input.InviteCode,
		}
		if err = tx.Where(&u).Find(&u).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				err = common_error.ErrInvalidInviteCode
			}
			return
		}

		inviter = &u
	}

	userInfo := user_model.User{
		Username: username,
		Nickname: &username,
		Password: util.GeneratePassword(input.Password),
		Status:   user_model.UserStatusInactivated, // 开始时未激活状态
		Role:     pq.StringArray{role_model.DefaultUser.Name},
		Phone:    input.Phone,
		Email:    input.Email,
		Gender:   user_model.GenderErrUnknown,
	}

	if err = tx.Create(&userInfo).Error; err != nil {
		return
	}

	// 如果存在邀请者的话，写入邀请列表中
	if inviter != nil {
		inviteHistory := invite_model.InviteHistory{
			Inviter:       inviter.Id,
			Invitee:       userInfo.Id,
			Status:        invite_model.StatusInviteRegistered,
			RewardSettled: false,
		}

		// 创建邀请记录
		if err = tx.Create(&inviteHistory).Error; err != nil {
			return
		}
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// 创建用户对应的钱包账号
	for _, walletName := range wallet_model.Wallets {
		if err = tx.Table(wallet.GetTableName(walletName)).Create(&wallet_model.Wallet{
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
		if err = redis.ActivationCodeClient.Set(activationCode, userInfo.Id, time.Minute*30).Err(); err != nil {
			return
		}

		// 把 "发送激活码" 加入消息队列
		var body []byte

		if body, err = json.Marshal(message_queue.SendActivationEmailBody{
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

func SignUpRouter(ctx *gin.Context) {
	var (
		input SignUpParams
		err   error
		res   = schema.Response{}
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

	res = SignUp(input)
}
