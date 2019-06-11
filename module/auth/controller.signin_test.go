// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/common_error"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestSignInWithEmptyBody(t *testing.T) {
	// empty body
	r := tester.HttpUser.Post("/v1/auth/signin", []byte(nil), nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, common_error.ErrInvalidParams.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestSignInWithErrInvalidPassword(t *testing.T) {
	body, _ := json.Marshal(&auth.SignInParams{
		Account:  "TestSignInWithErrInvalidPassword",
		Password: "abc", // 输入错误的密码
	})

	// empty body
	r := tester.HttpUser.Post("/v1/auth/signin", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, common_error.ErrInvalidAccountOrPassword.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestSignInSuccess(t *testing.T) {
	rand.Seed(111)
	// 先注册一个账号
	username := "test-TestSignInSuccess"
	password := "123123"

	if r := auth.SignUp(auth.SignUpParams{
		Username: &username,
		Password: password,
	}); r.Status != schema.StatusSuccess {
		t.Error(r.Message)
		return
	} else {
		defer func() {
			auth.DeleteUserByUserName(username)
		}()
	}

	res := auth.SignIn(schema.Context{
		UserAgent: "test-user-agent",
		Ip:        "0.0.0.0.0",
	}, auth.SignInParams{
		Account:  username,
		Password: password,
	})

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := user_schema.ProfileWithToken{}

	if err := mapstructure.Decode(res.Data, &profile); err != nil {
		assert.Error(t, err, err.Error())
	}

	assert.NotEmpty(t, profile.Token)

	if c, err := token.Parse("Bearer "+profile.Token, false); err != nil {
		t.Error(err)
		return
	} else {
		assert.IsType(t, "", c.Uid, "UID必须是字符串")
	}
}
