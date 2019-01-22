package user_test

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/src/util"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetProfileWithInvalidAuth(t *testing.T) {

	header := mocker.Header{
		"Authorization": "Bearera 12312", // invalid Bearera
	}

	r := tester.Http.Get("/v1/user/profile", []byte(""), &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}

	if !assert.Equal(t, schema.StatusFail, res.Status) {
		return
	}
	if !assert.Equal(t, exception.InvalidAuth.Error(), res.Message) {
		return
	}
}

func TestGetProfileWithInvalidToken(t *testing.T) {
	header := mocker.Header{
		"Authorization": util.TokenPrefix + " 12312",
	}

	r := tester.Http.Get("/v1/user/profile", []byte(""), &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}

	if !assert.Equal(t, schema.StatusFail, res.Status) {
		return
	}
	if !assert.Equal(t, exception.InvalidToken.Error(), res.Message) {
		return
	}
}

func TestGetProfile(t *testing.T) {
	var (
		uid         string
		username    = "test-TestResetPasswordSuccess"
		password    = "123123"
		tokenString string
	)
	if r := auth.SignUp(auth.SignUpParams{
		Username: &username,
		Password: password,
	}); r.Status != schema.StatusSuccess {
		t.Error(r.Message)
		return
	} else {
		userInfo := schema.Profile{}
		if err := tester.Decode(r.Data, &userInfo); err != nil {
			t.Error(err)
			return
		}
		uid = userInfo.Id
		defer func() {
			auth.DeleteUserByUserName(username)
		}()

		// 登陆获取Token
		if r := auth.SignIn(controller.Context{
			UserAgent: "test",
			Ip:        "0.0.0.0.0",
		}, auth.SignInParams{
			Account:  username,
			Password: password,
		}); r.Status != schema.StatusSuccess {
			t.Error(r.Message)
			return
		} else {
			userInfo := schema.ProfileWithToken{}
			if err := tester.Decode(r.Data, &userInfo); err != nil {
				t.Error(err)
				return
			}
			tokenString = userInfo.Token
		}
	}

	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + tokenString,
	}

	r := tester.Http.Get("/v1/user/profile", []byte(""), &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}

	if !assert.Equal(t, schema.StatusSuccess, res.Status) {
		fmt.Println(res.Message)
		return
	}
	if !assert.Equal(t, "", res.Message) {
		return
	}

	profile := schema.Profile{}

	if assert.Nil(t, tester.Decode(res.Data, &profile)) {
		return
	}

	if !assert.Equal(t, uid, profile.Id) {
		return
	}
	if !assert.Equal(t, username, *profile.Email) {
		return
	}
}
