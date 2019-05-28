package middleware

import (
	"errors"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/database"
	"github.com/axetroy/go-server/src/util"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

var (
	PayPasswordHeader = "X-Pay-Password"
)

// 交易密码的验证中间件
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
			context.JSON(http.StatusOK, schema.Response{
				Status:  schema.StatusFail,
				Message: err.Error(),
				Data:    nil,
			})

			// 中断后面的路由器执行
			context.Abort()

			return
		}
	}()

	payPassword := context.GetHeader(PayPasswordHeader)

	if len(payPassword) == 0 {
		err = exception.RequirePayPassword
		return
	}

	uid := context.GetString(ContextUidField)

	if uid == "" {
		err = exception.UserNotLogin
		return
	}

	userInfo := model.User{Id: uid}

	if err = database.Db.Where(&userInfo).Last(&userInfo).Error; err != nil {
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
	if *userInfo.PayPassword != util.GeneratePassword(payPassword) {
		err = exception.InvalidPassword
		return
	}

}
