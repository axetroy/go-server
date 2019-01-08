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
	var (
		username    = "test-TestResetPasswordSuccess"
		oldPassword = "123123"
		uid         string
		resetCode   string
		newPassword = "321321"
	)
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

	resetCode = email.GenerateResetCode(uid)

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

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, response.StatusSuccess, res.Status)
	assert.Equal(t, true, res.Data)

	var (
		ormErr error
	)

	tx := orm.DB.Begin()

	defer func() {
		// 重置密码回旧密码
		userInfo := &model.User{
			Id: uid,
		}

		ormErr = tx.Model(&userInfo).Update("password", password.Generate(oldPassword)).Error

		if ormErr == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	userInfo := &model.User{
		Id: uid,
	}

	if ormErr = tx.Where(&userInfo).First(&userInfo).Error; ormErr != nil {
		t.Error(ormErr)
		return
	}

	// 两次密码应该一致
	assert.Equal(t, password.Generate(newPassword), userInfo.Password)
}
