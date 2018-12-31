package email

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/email"
	"github.com/axetroy/go-server/services/redis"
	"net/http"
	"strconv"
	"time"
)

type SendActivationEmailParams struct {
	To string `json:"to"` // 发送给谁
}

func GenerateActivationCode(uid int64) string {
	// 生成重置码
	activationCode := "activation-" + strconv.FormatInt(uid, 10)
	return activationCode
}

func SendActivationEmail(context *gin.Context) {
	var (
		input   SendActivationEmailParams
		err     error
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
				Data:    true,
			})
		}
	}()

	if err = context.ShouldBindJSON(&input); err != nil {
		err = exception.InvalidParams
		return
	}

	// TODO: 校验参数

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	tx = true

	user := model.User{}

	var isExist bool

	if isExist, err = session.Where("email = ?", input.To).Get(&user); err != nil {
		return
	}

	if isExist == false {
		err = exception.UserNotExist
		return
	}

	if user.Status != model.UserStatusInactivated {
		err = exception.UserHaveActive
		return
	}

	// generate activation code
	activationCode := GenerateActivationCode(user.Id)

	// set activationCode to redis
	if err = redis.ActivationCode.Set(activationCode, user.Id, time.Minute*30).Err(); err != nil {
		return
	}

	e := email.New()

	// send email
	if err = e.SendActivationEmail(input.To, activationCode); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ActivationCode.Del(activationCode).Err()
		return
	}

}
