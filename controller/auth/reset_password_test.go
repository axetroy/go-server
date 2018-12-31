package auth_test

import (
	"encoding/json"
	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
	"gitlab.com/axetroy/server/controller/auth"
	"gitlab.com/axetroy/server/controller/email"
	"gitlab.com/axetroy/server/exception"
	"gitlab.com/axetroy/server/model"
	"gitlab.com/axetroy/server/orm"
	"gitlab.com/axetroy/server/response"
	"gitlab.com/axetroy/server/services/password"
	"gitlab.com/axetroy/server/services/redis"
	"gitlab.com/axetroy/server/tester"
	"net/http"
	"testing"
	"time"
)

func TestResetPasswordWithEmptyBody(t *testing.T) {
	t.Skip()
	// empty body
	r := tester.Http.Put("/v1/auth/password/reset", []byte(nil), nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusFail, res.Status)
	assert.Equal(t, exception.InvalidParams.Error(), res.Message, )
	assert.Nil(t, res.Data)
}

func TestResetPasswordWithInvalidPassword(t *testing.T) {
	t.Skip()
	newPassword := "321321"

	body, _ := json.Marshal(&auth.ResetPasswordParams{
		NewPassword: newPassword,
		Code:        "123123", // 错误的重置码
	})

	// empty body
	r := tester.Http.Put("/v1/auth/password/reset", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, res.Status, response.StatusFail)
	assert.Equal(t, res.Message, exception.InvalidResetCode.Error())
	assert.Nil(t, res.Data)
}

func TestResetPasswordSuccess(t *testing.T) {
	newPassword := "321321"

	var resetCode = email.GenerateResetCode(tester.Uid)

	// set to redis
	// set activationCode to redis
	if err := redis.ResetCode.Set(resetCode, tester.Uid, time.Minute*30).Err(); err != nil {
		t.Error(err)
		return
	}

	body, _ := json.Marshal(&auth.ResetPasswordParams{
		NewPassword: newPassword,
		Code:        resetCode,
	})

	// empty body
	r := tester.Http.Put("/v1/auth/password/reset", body, nil)

	if ok := assert.Equal(t, http.StatusOK, r.Code); !ok {
		return
	}

	res := response.Response{}

	if ok := assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)); !ok {
		return
	}

	if ok := assert.Equal(t, response.StatusSuccess, res.Status); !ok {
		return
	}
	if ok := assert.Equal(t, true, res.Data); !ok {
		return
	}

	// 查询密码是否已修改正确
	var (
		session *xorm.Session
		err     error
	)

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = session.Rollback()
		} else {
			_ = session.Commit()
		}
	}()

	defer func() {
		// 重置密码回旧密码
		user := &model.User{}
		var isExist bool
		if isExist, err = session.Where("email = ?", tester.Username).Get(user); err != nil {
			return
		}
		if isExist == false {
			err = exception.UserNotExist
			return
		}

		user.Password = password.Generate(tester.Password)

		if _, er := session.Where("email = ?", tester.Username).Cols("password").Update(user); er != nil {
			err = er
		}
	}()

	user := model.User{
		Email: &tester.Username,
	}

	if isExist, err := session.Get(&user); err != nil {
		t.Error(err)
		return
	} else {
		if isExist == false {
			err = exception.UserNotExist
			t.Error(err)
			return
		}
	}

	assert.Equal(t, password.Generate(newPassword), user.Password)
}
