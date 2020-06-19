package auth

import (
	"context"
	"encoding/json"
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
	"github.com/axetroy/go-server/internal/service/email"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/internal/service/telephone"
	"github.com/axetroy/go-server/pkg/proto"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"github.com/mssola/user_agent"
	"net/url"
	"time"
)

type SendEmailAuthCodeParams struct {
	Email string `json:"email" validate:"required,email" comment:"邮箱"`
}

type SendPhoneAuthCodeParams struct {
	Phone string `json:"phone" validate:"required,numeric,len=11" comment:"手机号"`
}

func GenerateAuthCode() string {
	return util.RandomString(6)
}

// 发送邮箱验证码 (不需要登陆)
func SendEmailAuthCode(c helper.Context, input SendEmailAuthCodeParams) (res schema.Response) {
	var (
		err error
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

		helper.Response(&res, nil, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	// 生成验证码
	activationCode := GenerateAuthCode()

	// 缓存验证码到 redis
	if err = redis.ClientAuthEmailCode.Set(context.Background(), activationCode, input.Email, time.Minute*10).Err(); err != nil {
		return
	}

	e, err := email.NewMailer()

	if err != nil {
		return
	}

	// send email
	if err = e.SendAuthEmail(input.Email, activationCode); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ClientAuthEmailCode.Del(context.Background(), activationCode).Err()
		return
	}

	return
}

// 发送手机验证码 (不需要登陆)
func SendPhoneAuthCode(c helper.Context, input SendPhoneAuthCodeParams) (res schema.Response) {
	var (
		err error
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

		helper.Response(&res, nil, nil, err)
	}()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	// 生成验证码
	activationCode := GenerateAuthCode()

	// 缓存验证码到 redis
	if err = redis.ClientAuthPhoneCode.Set(context.Background(), activationCode, input.Phone, time.Minute*10).Err(); err != nil {
		return
	}

	if err = telephone.GetClient().SendAuthCode(input.Phone, activationCode); err != nil {
		// 如果发送失败，则删除
		_ = redis.ClientAuthPhoneCode.Del(context.Background(), activationCode).Err()
		return
	}

	return
}

var SendEmailAuthCodeRouter = router.Handler(func(c router.Context) {
	var (
		input SendEmailAuthCodeParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SendEmailAuthCode(helper.NewContext(&c), input)
	})
})

var SendPhoneAuthCodeRouter = router.Handler(func(c router.Context) {
	var (
		input SendPhoneAuthCodeParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return SendPhoneAuthCode(helper.NewContext(&c), input)
	})
})

type QRCodeBody struct {
	SessionID string `json:"session_id"` // 会话 ID
	ExpiredAt string `json:"expired_at"` // 会话过期时间
}

type QRCodeEntry struct {
	OS      string  `json:"os"`                // 操作系统
	Browser string  `json:"browser"`           // 浏览器名称
	Version string  `json:"version"`           // 浏览器版本
	Ip      string  `json:"ip"`                // IP 地址
	UserID  *string `json:"user_id,omitempty"` // 对应的 user ID
}

func QRCodeGenerateLoginLink(c helper.Context) (res schema.Response) {
	var (
		err  error
		data string
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

		helper.Response(&res, data, nil, err)
	}()

	expAt := time.Minute * 5

	sessionID, err := uuid.NewRandom()

	if err != nil {
		return
	}

	body := QRCodeBody{
		SessionID: sessionID.String(),
		ExpiredAt: time.Now().Add(expAt).Format(time.RFC3339Nano),
	}

	link, err := proto.NewProto(proto.Auth, body).String()

	if err != nil {
		return
	}

	data = link

	ua := user_agent.New(c.UserAgent)

	browserName, browserVersion := ua.Browser()

	entry := QRCodeEntry{
		OS:      ua.OS(),
		Browser: browserName,
		Version: browserVersion,
		Ip:      c.Ip,
	}

	entryStr, err := json.Marshal(entry)

	if err != nil {
		return
	}

	if err = redis.QRCodeLoginCode.Set(context.Background(), sessionID.String(), entryStr, expAt).Err(); err != nil {
		return
	}

	return
}

var QRCodeGenerateLoginLinkRouter = router.Handler(func(c router.Context) {
	c.ResponseFunc(nil, func() schema.Response {
		return QRCodeGenerateLoginLink(helper.NewContext(&c))
	})
})

type QRCodeCheckParams struct {
	Url      string `json:"url" validate:"required" comment:"URL"`
	Duration *int64 `json:"duration" validate:"omitempty" comment:"持续时间"`
}

func QRCodeLoginCheck(c helper.Context, input QRCodeCheckParams) (res schema.Response) {
	var (
		err  error
		data = &schema.ProfileWithToken{}
		tx   *gorm.DB
		u    *url.URL
		p    *proto.Proto
		b    []byte
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

		if err != nil {
			data = nil
		}

		helper.Response(&res, data, nil, err)
	}()

	tx = database.Db.Begin()

	// 参数校验
	if err = validator.ValidateStruct(input); err != nil {
		return
	}

	u, err = url.Parse(input.Url)

	if err != nil {
		return
	}

	switch u.Scheme {
	case proto.Auth.String():
		{
			p, err = proto.Parse(input.Url)

			if err != nil {
				return
			}

			b, err = p.Data()

			if err != nil {
				return
			}

			var payload QRCodeBody

			if err = json.Unmarshal(b, &payload); err != nil {
				return
			}

			var val string

			val, err = redis.QRCodeLoginCode.Get(context.Background(), payload.SessionID).Result()

			if err != nil {
				return
			}

			var entry QRCodeEntry

			if err = json.Unmarshal([]byte(val), &entry); err != nil {
				return
			}

			if entry.UserID == nil {
				err = exception.NoData
				if err != nil {
					return
				}
				return
			}

			userInfo := model.User{}

			if err = tx.Model(userInfo).Where("id = ?", entry.UserID).First(&userInfo).Error; err != nil {
				return
			}

			if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
				return
			}

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
			if t, err := authentication.Gateway(false).Generate(userInfo.Id, duration); err != nil {
				return
			} else {
				data.Token = t
			}

			// 写入登陆记录
			loginLog := model.LoginLog{
				Uid:     userInfo.Id,                       // 用户ID
				Type:    model.LoginLogTypeQRCode,          // 默认用户名登陆
				Command: model.LoginLogCommandLoginSuccess, // 登陆成功
				Client:  c.UserAgent,                       // 用户的 userAgent
				LastIp:  c.Ip,                              // 用户的IP
			}

			if err = tx.Create(&loginLog).Error; err != nil {
				return
			}

			// 删除 redis
			if err = redis.QRCodeLoginCode.Del(context.Background(), payload.SessionID).Err(); err != nil {
				return
			}
		}
	default:
		err = exception.InvalidParams
		return
	}

	return
}

var QRCodeLoginCheckRouter = router.Handler(func(c router.Context) {
	var (
		input QRCodeCheckParams
	)
	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return QRCodeLoginCheck(helper.NewContext(&c), input)
	})
})
