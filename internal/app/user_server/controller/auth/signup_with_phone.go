package auth

import (
	"context"
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"time"
)

type SignUpWithPhoneParams struct {
	Phone      string  `json:"phone" validate:"required,numeric,len=11" comment:"手机号"` // 手机号
	Code       string  `json:"code" validate:"required" comment:"验证码"`                 // 短信验证码
	InviteCode *string `json:"invite_code" validate:"omitempty,len=8" comment:"邀请码"`   // 邀请码
}

// 使用手机注册
func SignUpWithPhone(input SignUpWithPhoneParams) (res schema.Response) {
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

	phone, err := redis.ClientAuthPhoneCode.Get(context.Background(), input.Code).Result()

	if err != nil {
		return
	}

	// 校验短信验证码是否一致
	if phone != input.Phone {
		err = exception.InvalidParams
		return
	}

	tx = database.Db.Begin()

	u := model.User{Phone: &input.Phone}

	if err = tx.Where("phone = ?", input.Phone).Find(&u).Error; err != nil {
		if err != gorm.ErrRecordNotFound {
			return
		}
	}

	if u.Id != "" {
		err = exception.UserExist
		return
	}

	username := "u" + util.GenerateId()
	pwd := util.RandomString(6)

	userInfo := model.User{
		Username:                username,
		Nickname:                &username,
		Password:                util.GeneratePassword(pwd),
		Status:                  model.UserStatusInit,
		Role:                    pq.StringArray{model.DefaultUser.Name},
		Phone:                   &input.Phone,
		Email:                   nil,
		Gender:                  model.GenderUnknown,
		UsernameRenameRemaining: 1, // 允许重命名 username 一次
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

var SignUpWithPhoneRouter = router.Handler(func(c router.Context) {
	var (
		input SignUpWithPhoneParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignUpWithPhone(input)
	})
})
