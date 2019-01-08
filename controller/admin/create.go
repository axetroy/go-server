package admin

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/password"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/mitchellh/mapstructure"
	"net/http"
	"time"
)

type CreateAdminParams struct {
	Account  string `json:"account"`  // 管理员账号，登陆凭借
	Password string `json:"password"` // 管理员密码
	Name     string `json:"name"`     // 管理员名称，注册后不可修改
}

// 创建管理员
func CreateAdmin(input CreateAdminParams, isSuper bool) (res response.Response) {
	var (
		err  error
		data Detail
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

	tx = orm.DB.Begin()

	n := model.Admin{Username: input.Account}

	if tx.Where(&n).First(&n).RecordNotFound() == false {
		err = exception.AdminExist
		return
	}

	adminInfo := model.Admin{
		Username: input.Account,
		Name:     input.Name,
		Password: password.Generate(input.Password),
		Status:   model.AdminStatusInit,
		IsSuper:  isSuper,
	}

	if err = tx.Create(&adminInfo).Error; err != nil {
		return
	}

	if err = mapstructure.Decode(adminInfo, &data.Pure); err != nil {
		return
	}

	data.CreatedAt = adminInfo.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminInfo.UpdatedAt.Format(time.RFC3339Nano)

	return
}

func CreateAdminRouter(context *gin.Context) {
	var (
		input CreateAdminParams
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

	uid := context.GetString("uid")

	adminInfo := model.Admin{
		Id: uid,
	}

	if err = orm.DB.Where(&adminInfo).First(&adminInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.AdminNotExist
			return
		}
		return
	}

	if adminInfo.IsSuper == false {
		err = exception.AdminNotSuper
		return
	}

	res = CreateAdmin(input, false)
}
