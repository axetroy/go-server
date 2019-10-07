// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/email"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/service/redis"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestActivationWithEmptyBody(t *testing.T) {

	// empty body
	r := tester.HttpUser.Post("/v1/auth/activation", []byte(nil), nil)

	if ok := assert.Equal(t, http.StatusOK, r.Code); !ok {
		return
	}

	res := schema.Response{}

	if ok := assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res)); !ok {
		return
	}

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, exception.InvalidParams.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestActivationWithInvalidCode(t *testing.T) {
	body, _ := json.Marshal(&auth.ActivationParams{
		Code: "123",
	})

	// empty body
	r := tester.HttpUser.Post("/v1/auth/activation", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	assert.Equal(t, exception.InvalidActiveCode.Code(), res.Status)
	assert.Equal(t, exception.InvalidActiveCode.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestActivationSuccess(t *testing.T) {
	var (
		testerUsername = "tester-TestActivationSuccess"
		testerUid      = ""
	)
	// 动态创建一个测试账号
	{
		r := auth.SignUp(auth.SignUpParams{
			Username: &testerUsername,
			Password: "123123",
		}, model.UserStatusInactivated)

		profile := schema.Profile{}

		assert.Nil(t, tester.Decode(r.Data, &profile))

		testerUid = profile.Id

		defer func() {
			auth.DeleteUserByUserName(testerUsername)
		}()
	}

	// generate activation code
	activationCode := email.GenerateActivationCode(testerUid)

	// set activationCode to redis
	if err := redis.ClientActivationCode.Set(activationCode, testerUid, time.Minute*30).Err(); err != nil {
		t.Error(err)
		return
	}

	defer func() {
		// remove activation code
		_ = redis.ClientActivationCode.Del(activationCode).Err()
	}()

	body, _ := json.Marshal(&auth.ActivationParams{
		Code: activationCode,
	})

	// empty body
	r := tester.HttpUser.Post("/v1/auth/activation", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)
}
