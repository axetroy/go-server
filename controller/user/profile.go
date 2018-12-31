package user

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/mitchellh/mapstructure"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"net/http"
	"strconv"
	"time"
)

type ProfilePure struct {
	Id       int64            `json:"id"`
	Username string           `json:"username"`
	Nickname *string          `json:"nickname"`
	Email    *string          `json:"email"`
	Phone    *string          `json:"phone"`
	Status   model.UserStatus `json:"status"`
	Avatar   string           `json:"avatar"`
	Role     string           `json:"role"`
	Level    int32            `json:"level"`
}

type Profile struct {
	ProfilePure
	PayPassword bool   `json:"pay_password"` // 是否已设置交易密码
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateProfileParams struct {
	Nickname *string `json:"nickname"`
	Avatar   *string `json:"avatar"`
}

func GetProfile(context *gin.Context) {
	var (
		err     error
		uid     int64
		data    Profile
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

	if val, isExist := context.Get("uid"); isExist != true {
		return
	} else {
		if uid, err = strconv.ParseInt(fmt.Sprintf("%v", val), 10, 64); err != nil {
			return
		}
	}

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	user := model.User{Id: uid}

	var isExist bool

	if isExist, err = session.Get(&user); err != nil {
		return
	}

	if isExist == false {
		err = exception.UserNotExist
		return
	}

	if err = mapstructure.Decode(user, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = user.PayPassword != nil && len(*user.PayPassword) != 0
	data.CreatedAt = user.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = user.UpdatedAt.Format(time.RFC3339Nano)
}

func UpdateProfile(context *gin.Context) {
	var (
		input   UpdateProfileParams
		err     error
		uid     int64
		data    Profile
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

	if val, isExist := context.Get("uid"); isExist != true {
		return
	} else {
		if uid, err = strconv.ParseInt(fmt.Sprintf("%v", val), 10, 64); err != nil {
			return
		}
	}

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	user := model.User{}

	query := session.Where("id = ?", uid)

	if isExist, er := query.Get(&user); er != nil {
		err = er
	} else {
		if isExist == false {
			err = exception.UserNotExist
			return
		}
	}

	if input.Nickname != nil {
		user.Nickname = input.Nickname
		query = query.Cols("nickname")
	}

	if input.Avatar != nil {
		user.Avatar = *input.Avatar
		query = query.Cols("avatar")
	}

	if _, err = query.Update(&user); err != nil {
		return
	}

	if err = mapstructure.Decode(user, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = user.PayPassword != nil && len(*user.PayPassword) != 0
	data.CreatedAt = user.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = user.UpdatedAt.Format(time.RFC3339Nano)
}
