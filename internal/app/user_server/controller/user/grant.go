package user

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/axetroy/go-server/internal/app/user_server/controller/auth"
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/library/router"
	"github.com/axetroy/go-server/internal/library/validator"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/pkg/proto"
	"net/url"
	"time"
)

type QRCodeAuthParams struct {
	Url string `json:"url" validate:"required" comment:"URL"`
}

func QRCodeAuthGrant(c helper.Context, input QRCodeAuthParams) (res schema.Response) {
	var (
		err error
		u   *url.URL
		p   *proto.Proto
		b   []byte
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

	u, err = url.Parse(input.Url)

	if err != nil {
		return
	}

	switch u.Scheme {
	case string(proto.Auth):
		{
			p, err = proto.Parse(input.Url)

			if err != nil {
				return
			}

			b, err = p.Data()

			if err != nil {
				return
			}

			var payload auth.QRCodeBody

			if err := json.Unmarshal(b, &payload); err != nil {
				return
			}

			val, err := redis.QRCodeLoginCode.Get(context.Background(), payload.SessionID).Result()

			if err != nil {
				return
			}

			var entry auth.QRCodeEntry

			if err := json.Unmarshal([]byte(val), &entry); err != nil {
				return
			}

			entry.UserID = &c.Uid

			b, err := json.Marshal(entry)

			if err != nil {
				return
			}

			// 更新 redis
			if err = redis.QRCodeLoginCode.Set(context.Background(), payload.SessionID, string(b), time.Minute*2).Err(); err != nil {
				return
			}
		}
	default:
		err = exception.InvalidParams
		return
	}

	return
}

var QRCodeAuthGrantRouter = router.Handler(func(c router.Context) {
	var (
		input QRCodeAuthParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return QRCodeAuthGrant(helper.NewContext(&c), input)
	})
})

func QRCodeAuthQuery(c helper.Context, link string) (res schema.Response) {
	var (
		err  error
		data auth.QRCodeEntry
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

	u, err := url.Parse(link)

	if err != nil {
		return
	}

	switch u.Scheme {
	case string(proto.Auth):
		{
			encodedStr := u.RawPath

			if b, err := base64.StdEncoding.DecodeString(encodedStr); err != nil {
				return
			} else {
				var payload auth.QRCodeBody

				if err := json.Unmarshal(b, &payload); err != nil {
					return
				}

				val, err := redis.QRCodeLoginCode.Get(context.Background(), payload.SessionID).Result()

				if err != nil {
					return
				}

				if err := json.Unmarshal([]byte(val), &data); err != nil {
					return
				}
			}
		}
	default:
		err = exception.InvalidParams
		return
	}

	return
}

var QRCodeAuthQueryRouter = router.Handler(func(c router.Context) {
	var (
		input QRCodeAuthParams
	)

	c.ResponseFunc(c.ShouldBindJSON(&input), func() schema.Response {
		return QRCodeAuthQuery(helper.NewContext(&c), c.Param("link"))
	})
})
