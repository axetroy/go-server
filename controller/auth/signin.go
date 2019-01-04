package auth

import (
	"errors"
	"github.com/asaskevich/govalidator"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/id"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/password"
	"github.com/axetroy/go-server/token"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type SignInParams struct {
	Account  string  `json:"account"`
	Password string  `json:"password"`
	Code     *string `json:"code"` // 手机验证码
}

type SignInContext struct {
	UserAgent string
	Ip        string
}

type SignInResponse struct {
	user.Profile
	Token string `json:"token"`
}

func SignIn(input SignInParams, context SignInContext) (res response.Response) {
	var (
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
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = response.StatusSuccess
		}

	}()

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	userInfo := model.User{Password: password.Generate(input.Password)}

	if govalidator.Matches(input.Account, "^/d+$") && input.Code != nil { // 如果是手机号, 并且传入了code字段
		// TODO: 这里手机登陆应该用验证码作为校验
		userInfo.Phone = &input.Account
	} else if govalidator.IsEmail(input.Account) { // 如果是邮箱的话
		userInfo.Email = &input.Account
	} else {
		userInfo.Username = input.Account // 其他则为用户名
	}

	var hasExist bool

	if hasExist, err = session.Get(&userInfo); err != nil {
		return
	}

	if hasExist == false {
		err = exception.InvalidAccountOrPassword
		return
	}

	if err = mapstructure.Decode(userInfo, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = userInfo.PayPassword != nil && len(*userInfo.PayPassword) != 0
	data.CreatedAt = userInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = userInfo.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	if t, er := token.Generate(userInfo.Id, false); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	// 写入登陆记录
	var log = &model.LoginLog{
		Id:       id.Generate(),
		Uid:      userInfo.Id,
		Username: userInfo.Username,
		Type:     0, // 默认用户名登陆
		Command:  1, // 登陆成功
		Client:   context.UserAgent,
		LastIp:   context.Ip,
	}

	if _, err = session.Insert(log); err != nil {
		return
	}

	return
}

func SignInRouter(context *gin.Context) {
	var (
		input SignInParams
		err   error
		res   = response.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	res = SignIn(input, SignInContext{
		UserAgent: context.GetHeader("user-agent"),
		Ip:        context.ClientIP(),
	})
}
