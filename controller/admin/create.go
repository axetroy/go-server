package admin

import (
	"errors"
	"fmt"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/id"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
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
		err     error
		data    Detail
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

	if hasExist, er := session.Exist(&model.Admin{
		Username: input.Account,
	}); er != nil {
		err = er
		return
	} else if hasExist {
		err = exception.AdminExist
		return
	}

	adminInfo := model.Admin{
		Id:       id.Generate(),
		Username: input.Account,
		Name:     input.Name,
		Password: input.Password,
		Status:   model.AdminStatusInit,
		IsSuper:  isSuper,
	}

	if _, err = session.Insert(&adminInfo); err != nil {
		return
	}

	adminData := model.Admin{
		Id: adminInfo.Id,
	}

	if isExist, er := session.Get(&adminData); er != nil {
		err = er
		return
	} else if !isExist {
		err = exception.New("创建失败")
		return
	}

	if err = mapstructure.Decode(adminData, &data.Pure); err != nil {
		return
	}

	data.CreatedAt = adminData.CreatedAt.Format(time.RFC3339Nano)
	data.UpdatedAt = adminData.UpdatedAt.Format(time.RFC3339Nano)

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

	// TODO: 验证是否是超级管理员
	uid := context.GetString("uid")

	fmt.Println(uid)

	res = CreateAdmin(input, false)
}
