package auth_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/controller/email"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/password"
	"github.com/axetroy/go-server/services/redis"
	"github.com/axetroy/go-server/tester"
	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestResetPasswordWithEmptyBody(t *testing.T) {
	// empty body
	r := tester.Http.Put("/v1/auth/password/reset", []byte(nil), nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusFail, res.Status)
	assert.Equal(t, exception.InvalidParams.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestResetPasswordWithInvalidPassword(t *testing.T) {
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
	assert.Equal(t, false, res.Data)
}

func TestResetPasswordSuccess(t *testing.T) {
	// 先创建一个测试账号
	username := "test-TestResetPasswordSuccess"
	oldPassword := "123123"
	var uid string
	if r := auth.SignUp(auth.SignUpParams{
		Username: &username,
		Password: oldPassword,
	}); r.Status != response.StatusSuccess {
		t.Error(r.Message)
		return
	} else {
		userInfo := user.Profile{}
		if err := tester.Decode(r.Data, &userInfo); err != nil {
			t.Error(err)
			return
		}
		uid = userInfo.Id
		defer func() {
			auth.DeleteUserByUserName(username)
		}()
	}

	newPassword := "321321"

	var resetCode = email.GenerateResetCode(uid)

	// set to redis
	// set activationCode to redis
	if err := redis.ResetCode.Set(resetCode, uid, time.Minute*30).Err(); err != nil {
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

	if !assert.Equal(t, response.StatusSuccess, res.Status) {
		return
	}

	if !assert.Equal(t, true, res.Data) {
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
		if isExist, err = session.Where("email = ?", username).Get(user); err != nil {
			return
		}
		if isExist == false {
			err = exception.UserNotExist
			return
		}

		user.Password = password.Generate(oldPassword)

		if _, er := session.Where("email = ?", username).Cols("password").Update(user); er != nil {
			err = er
		}
	}()

	user := model.User{
		Username: username,
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
