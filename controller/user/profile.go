package user

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type ProfilePure struct {
	Id         string           `json:"id"`
	Username   string           `json:"username"`
	Nickname   *string          `json:"nickname"`
	Email      *string          `json:"email"`
	Phone      *string          `json:"phone"`
	Status     model.UserStatus `json:"status"`
	Avatar     string           `json:"avatar"`
	Role       string           `json:"role"`
	Level      int32            `json:"level"`
	InviteCode string           `json:"invite_code"`
}

type Profile struct {
	ProfilePure
	PayPassword bool   `json:"pay_password"` // 是否已设置交易密码
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

type UpdateProfileParams struct {
	Nickname *string       `json:"nickname" valid:"length(1|36)~昵称长度为1-36位"`
	Gender   *model.Gender `json:"gender"`
	Avatar   *string       `json:"avatar"`
}

func GetProfile(uid string) (res response.Response) {
	var (
		err  error
		data Profile
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
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = response.StatusSuccess
		}
	}()

	tx = orm.DB.Begin()

	user := model.User{Id: uid}

	if err = tx.Where(&user).Last(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if err = mapstructure.Decode(user, &data.ProfilePure); err != nil {
		return
	}

	data.PayPassword = user.PayPassword != nil && len(*user.PayPassword) != 0
	data.CreatedAt = user.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = user.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func GetProfileRouter(context *gin.Context) {
	var (
		err error
		res = response.Response{}
	)

	defer func() {
		if err != nil {
			res.Data = nil
			res.Message = err.Error()
		}
		context.JSON(http.StatusOK, res)
	}()

	res = GetProfile(context.GetString("uid"))
}

func UpdateProfile(uid string, input UpdateProfileParams) (res response.Response) {
	var (
		err  error
		data Profile
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
			res.Message = err.Error()
			res.Data = nil
		} else {
			res.Data = data
			res.Status = response.StatusSuccess
		}
	}()

	// TODO: 参数校验

	tx = orm.DB.Begin()

	userInfo := model.User{
		Id: uid,
	}

	if err = tx.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	updateMap := map[string]interface{}{}

	if input.Nickname != nil {
		updateMap["nickname"] = input.Nickname
	}

	if input.Avatar != nil {
		updateMap["avatar"] = *input.Avatar
	}

	if input.Gender != nil {
		updateMap["gender"] = *input.Gender
	}

	if err = tx.Model(&userInfo).Updates(updateMap).Error; err != nil {
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

func UpdateProfileRouter(context *gin.Context) {
	var (
		err   error
		res   = response.Response{}
		input UpdateProfileParams
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

	res = UpdateProfile(context.GetString("uid"), input)
}
