package auth_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/token"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestSignInWithEmptyBody(t *testing.T) {
	// empty body
	r := tester.Http.Post("/v1/auth/signin", []byte(nil), nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusFail, res.Status)
	assert.Equal(t, exception.InvalidParams.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestSignInWithInvalidPassword(t *testing.T) {
	body, _ := json.Marshal(&auth.SignInParams{
		Account:  "TestSignInWithInvalidPassword",
		Password: "abc", // 输入错误的密码
	})

	// empty body
	r := tester.Http.Post("/v1/auth/signin", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, res.Status, response.StatusFail)
	assert.Equal(t, res.Message, exception.InvalidAccountOrPassword.Error())
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
	}); r.Status != response.StatusSuccess {
		t.Error(r.Message)
		return
	} else {
		defer func() {
			auth.DeleteUserByUserName(username)
		}()
	}

	res := auth.SignIn(auth.SignInParams{
		Account:  username,
		Password: password,
	}, auth.SignInContext{
		UserAgent: "test-user-agent",
		Ip:        "0.0.0.0.0",
	})

	assert.Equal(t, response.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := auth.SignInResponse{}

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
