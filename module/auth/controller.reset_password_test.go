// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/email"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_model"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/database"
	"github.com/axetroy/go-server/service/redis"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/util"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestResetPasswordWithEmptyBody(t *testing.T) {
	// empty body
	r := tester.HttpUser.Put("/v1/auth/password/reset", []byte(nil), nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, exception.ErrInvalidParams.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestResetPasswordWithErrInvalidPassword(t *testing.T) {
	newPassword := "321321"

	body, _ := json.Marshal(&auth.ResetPasswordParams{
		NewPassword: newPassword,
		Code:        "123123", // 错误的重置码
	})

	// empty body
	r := tester.HttpUser.Put("/v1/auth/password/reset", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, user_error.ErrInvalidResetCode.Error(), res.Message)
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
	}); r.Status != schema.StatusSuccess {
		t.Error(r.Message)
		return
	} else {
		userInfo := user_schema.Profile{}
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
	if err := redis.ResetCodeClient.Set(resetCode, uid, time.Minute*30).Err(); err != nil {
		t.Error(err)
		return
	}

	body, _ := json.Marshal(&auth.ResetPasswordParams{
		NewPassword: newPassword,
		Code:        resetCode,
	})

	// empty body
	r := tester.HttpUser.Put("/v1/auth/password/reset", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, true, res.Data)

	var (
		ormErr error
	)

	tx := database.Db.Begin()

	defer func() {
		// 重置密码回旧密码
		userInfo := &user_model.User{
			Id: uid,
		}

		ormErr = tx.Model(&userInfo).Update("password", util.GeneratePassword(oldPassword)).Error

		if ormErr == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	}()

	userInfo := &user_model.User{
		Id: uid,
	}

	if ormErr = tx.Where(&userInfo).First(&userInfo).Error; ormErr != nil {
		t.Error(ormErr)
		return
	}

	// 两次密码应该一致
	assert.Equal(t, util.GeneratePassword(newPassword), userInfo.Password)
}
