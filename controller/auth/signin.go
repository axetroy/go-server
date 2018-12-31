package auth

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/mitchellh/mapstructure"
	"gitlab.com/axetroy/server/controller/user"
	"gitlab.com/axetroy/server/exception"
	"gitlab.com/axetroy/server/id"
	"gitlab.com/axetroy/server/model"
	"gitlab.com/axetroy/server/orm"
	"gitlab.com/axetroy/server/response"
	"gitlab.com/axetroy/server/services/password"
	"gitlab.com/axetroy/server/token"
	"net/http"
	"time"
)

type SignInParams struct {
	Account  string  `json:"account"`
	Password string  `json:"password"`
	Code     *string `json:"code"` // 手机验证码
}

type SignInResponse struct {
	user.Profile
	Token string `json:"token"`
}

func Signin(context *gin.Context) {
	var (
		input   SignInParams
		err     error
		data    = &SignInResponse{}
		session *xorm.Session
		tx      bool
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
			session.Close()
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

	// TODO 校验input是否正确

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	user := model.User{Password: password.Generate(input.Password)}

	if govalidator.Matches(input.Account, "^/d+$") && input.Code != nil { // 如果是手机号, 并且传入了code字段
		// TODO: 这里手机登陆应该用验证码作为校验
		user.Phone = &input.Account
	} else if govalidator.IsEmail(input.Account) { // 如果是邮箱的话
		user.Email = &input.Account
	} else {
		user.Username = input.Account // 其他则为用户名
	}

	var hasExist bool

	if hasExist, err = session.Get(&user); err != nil {
		return
	}

	if hasExist == false {
		err = exception.InvalidAccountOrPassword
		return
	}

	if err = mapstructure.Decode(user, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = user.PayPassword != nil && len(*user.PayPassword) != 0
	data.CreatedAt = user.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = user.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	var tokenString string
	if tokenString, err = token.Generate(user.Id); err != nil {
		return
	}

	data.Token = tokenString

	// 写入登陆记录
	var log = &model.LoginLog{
		Id:       id.Generate(),
		Uid:      string(user.Id),
		Username: user.Username,
		Type:     0, // 默认用户名登陆
		Command:  1, // 登陆成功
		Client:   context.GetHeader("user-agent"),
		LastIp:   context.ClientIP(),
	}

	if _, err = session.Insert(log); err != nil {
		return
	}
}
