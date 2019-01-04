package admin

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/exception"
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
	Username string
	Password string
}

type ProfilePure struct {
	Id       string            `json:"id"`       // 用户ID
	Username string            `json:"username"` // 用户名, 用于登陆
	Name     string            `json:"name"`     // 管理员名
	IsSuper  bool              `json:"is_super"` // 是否是超级管理员, 超级管理员全站应该只有一个
	Status   model.AdminStatus `json:"status"`   // 状态
}

type Profile struct {
	ProfilePure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

type SignInResponse struct {
	Profile
	Token string `json:"token"`
}

func Login(input SignInParams) (res response.Response) {
	var (
		err     error
		data    = SignInResponse{}
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

	adminInfo := model.Admin{
		Username: input.Username,
		Password: password.Generate(input.Password),
	}

	var hasExist bool

	if hasExist, err = session.Get(&adminInfo); err != nil {
		return
	}

	if hasExist == false {
		err = exception.InvalidAccountOrPassword
		return
	}

	if err = mapstructure.Decode(adminInfo, &data.ProfilePure); err != nil {
		return
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	// generate token
	if t, er := token.Generate(adminInfo.Id, true); er != nil {
		err = er
		return
	} else {
		data.Token = t
	}

	return
}

func LoginRouter(context *gin.Context) {
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

	res = Login(input)

	fmt.Sprintf("%v", res)
}
