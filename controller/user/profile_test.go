package user_test

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/token"
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

	res := response.Response{}

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}

	if !assert.Equal(t, response.StatusFail, res.Status) {
		return
	}
	if !assert.Equal(t, exception.InvalidAuth.Error(), res.Message) {
		return
	}
}

func TestGetProfileWithInvalidToken(t *testing.T) {
	header := mocker.Header{
		"Authorization": token.Prefix + " 12312",
	}

	r := tester.Http.Get("/v1/user/profile", []byte(""), &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := response.Response{}

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}

	if !assert.Equal(t, response.StatusFail, res.Status) {
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

		// 登陆获取Token
		if r := auth.SignIn(auth.SignInParams{
			Account:  username,
			Password: password,
		}, auth.SignInContext{
			UserAgent: "test",
			Ip:        "0.0.0.0.0",
		}); r.Status != response.StatusSuccess {
			t.Error(r.Message)
			return
		} else {
			userInfo := auth.SignInResponse{}
			if err := tester.Decode(r.Data, &userInfo); err != nil {
				t.Error(err)
				return
			}
			tokenString = userInfo.Token
		}
	}

	header := mocker.Header{
		"Authorization": token.Prefix + " " + tokenString,
	}

	r := tester.Http.Get("/v1/user/profile", []byte(""), &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := response.Response{}

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}

	if !assert.Equal(t, response.StatusSuccess, res.Status) {
		fmt.Println(res.Message)
		return
	}
	if !assert.Equal(t, "", res.Message) {
		return
	}

	profile := user.Profile{}

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
