package auth

import (
	"errors"
	"github.com/axetroy/go-server/controller/invite"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/id"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/email"
	"github.com/axetroy/go-server/services/redis"
	"github.com/axetroy/redpack/services/password"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"strconv"
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

func SignUp(context *gin.Context) {
	var (
		input   SignUpParams
		err     error
		data    user.Profile
		session *xorm.Session
		tx      bool
		invitor *model.User
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

		if tx {
			if err != nil {
				_ = session.Rollback()
			} else {
				err = session.Commit()
			}
		}

		if session != nil {
			session.Close()
		}

		if err != nil {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})
		} else {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusSuccess,
				Message: "",
				Data:    data,
			})
		}

	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	if input.Password == "" {
		err = exception.RequirePassword
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
	}

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	var (
		username string
		uid      = id.Generate()
	)

	if input.Username == nil {
		username = "用户" + strconv.FormatInt(uid, 10)
	} else {
		username = *input.Username
	}

	var isExist bool

	if input.Username != nil {
		isExist, err = session.Where("username = ?", input.Username).Get(new(model.User))

		if err != nil {
			return
		}

		if isExist {
			err = exception.UserExist
			return
		}
	}

	if input.Email != nil {
		isExist, err = session.Where("email = ?", *input.Email).Get(new(model.User))

		if err != nil {
			return
		}

		if isExist {
			err = exception.UserExist
			return
		}
	}

	if input.Phone != nil {
		isExist, err = session.Where("phone = ?", *input.Phone).Get(new(model.User))

		if err != nil {
			return
		}

		if isExist {
			err = exception.UserExist
			return
		}
	}

	// 填入了邀请码，则去校验邀请码是否正确
	if input.InviteCode != nil {
		session := orm.Db.NewSession()
		u := model.User{
			InviteCode: *input.InviteCode,
		}
		if exist, er := session.Get(&u); er != nil {
			err = er
			return
		} else {
			if !exist {
				err = errors.New("无效的邀请码")
				return
			}
			invitor = &u
		}
	}

	userInfo := model.User{
		Id:         id.Generate(),
		Username:   username,
		Nickname:   &username,
		Password:   password.Generate(input.Password),
		Status:     model.UserStatusInactivated, // 开始时未激活状态
		Phone:      input.Phone,
		Email:      input.Email,
		InviteCode: invite.GenerateInviteCode(),
		Gender:     model.GenderUnknown,
	}

	if _, err = session.Insert(&userInfo); err != nil {
		return
	}

	if _, err = session.Get(&userInfo); err != nil {
		return
	}

	// 如果存在邀请者的话，写入邀请列表中
	if invitor != nil {
		inviteHistory := model.InviteHistory{
			Id:            id.Generate(),
			Invitor:       invitor.Id,
			Invited:       userInfo.Id,
			Status:        model.StatusInviteRegistered,
			RewardSettled: false,
		}

		if _, err = session.Insert(&inviteHistory); err != nil {
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
	cny := model.WalletCny{
		Wallet: model.Wallet{
			Id:      userInfo.Id,
			Balance: 0,
			Frozen:  0,
		},
	}
	usd := model.WalletUsd{
		Wallet: model.Wallet{
			Id:      userInfo.Id,
			Balance: 0,
			Frozen:  0,
		},
	}
	coin := model.WalletCoin{
		Wallet: model.Wallet{
			Id:      userInfo.Id,
			Balance: 0,
			Frozen:  0,
		},
	}
	if _, err = session.Insert(&cny); err != nil {
		return
	}
	if _, err = session.Insert(&usd); err != nil {
		return
	}
	if _, err = session.Insert(&coin); err != nil {
		return
	}

	// 如果是以邮箱注册的，那么发送激活链接
	if userInfo.Email != nil && len(*userInfo.Email) != 0 {
		// generate activation code
		activationCode := "activation-" + strconv.FormatInt(userInfo.Id, 10)

		// set activationCode to redis
		if err = redis.ActivationCode.Set(activationCode, userInfo.Id, time.Minute*30).Err(); err != nil {
			return
		}

		// send email
		mailer := email.New()

		// TODO: 把这个激活码放进队列, 因为发送邮箱实在是太慢了
		if err = mailer.SendActivationEmail(*input.Email, activationCode); err != nil {
			// 邮件没发出去的话，删除redis的key
			_ = redis.ActivationCode.Del(activationCode).Err()
			return
		}
		return
	}

}
