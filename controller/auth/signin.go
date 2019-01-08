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
	"github.com/jinzhu/gorm"
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
		err  error
		data = &SignInResponse{}
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

		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		} else {
			res.Data = data
			res.Status = response.StatusSuccess
		}

	}()

	userInfo := model.User{Password: password.Generate(input.Password)}

	if govalidator.Matches(input.Account, "^/d+$") && input.Code != nil { // 如果是手机号, 并且传入了code字段
		// TODO: 这里手机登陆应该用验证码作为校验
		userInfo.Phone = &input.Account
	} else if govalidator.IsEmail(input.Account) { // 如果是邮箱的话
		userInfo.Email = &input.Account
	} else {
		userInfo.Username = input.Account // 其他则为用户名
	}

	tx = orm.DB.Begin()

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.InvalidAccountOrPassword
		}
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
	log := model.LoginLog{
		Id:      id.Generate(),
		Uid:     userInfo.Id,
		Type:    model.LoginLogTypeUserName,        // 默认用户名登陆
		Command: model.LoginLogCommandLoginSuccess, // 登陆成功
		Client:  context.UserAgent,
		LastIp:  context.Ip,
	}

	if err = tx.Create(&log).Error; err != nil {
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
