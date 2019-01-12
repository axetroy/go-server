package middleware

import (
	"errors"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/password"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

// 验证交易密码的中间价
func AuthPayPassword(context *gin.Context) {
	var (
		err error
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

		// 如果有报错的话，那么不会进入到路由中
		if err != nil {
			context.JSON(http.StatusOK, response.Response{
				Status:  response.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})

			// 中断后面的路由器执行
			context.Abort()

			return
		}
	}()

	payPassword := context.GetHeader("X-Pay-Password")

	if len(payPassword) == 0 {
		err = exception.RequirePayPassword
		return
	}

	uid := context.GetString("uid")

	if uid == "" {
		err = exception.UserNotLogin
		return
	}

	userInfo := model.User{Id: uid}

	if err = orm.DB.Where(&userInfo).Last(&userInfo).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			err = exception.UserNotExist
		}
		return
	}

	if userInfo.PayPassword == nil {
		err = exception.PayPasswordNotSet
		return
	}

	// 校验密码是否正确
	if *userInfo.PayPassword != password.Generate(payPassword) {
		err = exception.InvalidPassword
		return
	}

}
