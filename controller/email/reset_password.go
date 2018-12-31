package email

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"gitlab.com/axetroy/server/exception"
	"gitlab.com/axetroy/server/id"
	"gitlab.com/axetroy/server/model"
	"gitlab.com/axetroy/server/orm"
	"gitlab.com/axetroy/server/response"
	"gitlab.com/axetroy/server/services/email"
	"gitlab.com/axetroy/server/services/redis"
	"gitlab.com/axetroy/server/utils"
	"net/http"
	"strconv"
	"time"
)

type SendResetPasswordEmailParams struct {
	To string `json:"to"` // 发送给谁
}

func GenerateResetCode(uid int64) string {
	// 生成重置码
	var codeId = "reset-" + strconv.FormatInt(id.Generate(), 10) + strconv.FormatInt(uid, 10)
	return utils.MD5(codeId)
}

func SendResetPasswordEmail(context *gin.Context) {
	var (
		input   SendResetPasswordEmailParams
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

	// 生成重置码
	var code = GenerateResetCode(user.Id)

	// set activationCode to redis
	if err = redis.ResetCode.Set(code, user.Id, time.Minute*30).Err(); err != nil {
		return
	}

	e := email.New()

	// send email
	if err = e.SendForgotPasswordEmail(input.To, code); err != nil {
		// 邮件没发出去的话，删除redis的key
		_ = redis.ResetCode.Del(code).Err()
		return
	}

}
