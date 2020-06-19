package auth

import (
	"errors"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/model"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/authentication"
	"github.com/axetroy/go-server/internal/service/database"
	"github.com/axetroy/go-server/internal/service/dotenv"
	"github.com/axetroy/go-server/internal/service/wechat"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"github.com/mitchellh/mapstructure"
	"time"
)

type SignInWithWechatParams struct {
	Code     string `json:"code" validate:"required" comment:"微信授权代码"` // 微信小程序授权之后返回的 code
	Duration *int64 `json:"duration" validate:"omitempty,number,gt=0" comment:"有效时间"`
}

// 使用微信小程序登陆
func SignInWithWechat(c helper.Context, input SignInWithWechatParams) (res schema.Response) {
	var (
		err  error
		data = &schema.ProfileWithToken{}
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

	wechatInfo, wechatErr := wechat.FetchOpenID(input.Code)

	if wechatErr != nil {
		err = wechatErr
		return
	}

	if len(wechatInfo.OpenID) == 0 {
		// 为了测试方便，那么我们就在测试环境下赋值一个假 ID， 否则会报错
		if dotenv.Test {
			wechatInfo.OpenID = "test open id"
		} else {
			err = exception.NoData
			return
		}

	}

	tx = database.Db.Begin()

	wechatOpenID := model.WechatOpenID{
		Id: wechatInfo.OpenID,
	}

	// 去查表
	result := tx.Where(&wechatOpenID).Preload("User").First(&wechatOpenID)

	var userInfo *model.User

	if result.RecordNotFound() {
		var (
			uid      = util.GenerateId()
			username = "v" + uid
		)

		userInfo = &model.User{
			Username:                username,
			Nickname:                &username,
			Password:                util.GeneratePassword(uid),
			Status:                  model.UserStatusInit,
			Role:                    pq.StringArray{model.DefaultUser.Name},
			Gender:                  model.GenderUnknown,
			WechatOpenID:            &wechatOpenID.Id,
			UsernameRenameRemaining: 1, // 允许微信注册的用户可以重命名一次
		}

		if err = CreateUserTx(tx, userInfo, nil); err != nil {
			return
		}

	} else {
		userInfo = &wechatOpenID.User
	}

	if userInfo == nil {
		err = exception.NoData
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	wechatBindingInfo := schema.WechatBindingInfo{}

	if err = mapstructure.Decode(wechatOpenID, &wechatBindingInfo); err != nil {
		return
	}

	data.Wechat = &wechatBindingInfo
	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	var duration time.Duration

	if input.Duration != nil {
		duration = time.Duration(*input.Duration * int64(time.Second))
	} else {
		duration = time.Hour * 6
	}

	// generate token
	if t, er := authentication.Gateway(false).Generate(userInfo.Id, duration); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	loginLog := model.LoginLog{
		Uid:     userInfo.Id,                       // 用户ID
		Type:    model.LoginLogTypeWechat,          // 微信登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  c.UserAgent,                       // 用户的 userAgent
		LastIp:  c.Ip,                              // 用户的IP
	}

	if err = tx.Create(&loginLog).Error; err != nil {
		return
	}

	return
}

var SignInWithWechatRouter = router.Handler(func(c router.Context) {
	var (
		input SignInWithWechatParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SignInWithWechat(helper.NewContext(&c), input)
	})
})
