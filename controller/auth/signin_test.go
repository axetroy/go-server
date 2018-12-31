package auth_test

import (
	"encoding/json"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/go-server/token"
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
	assert.Equal(t, exception.InvalidParams.Error(), res.Message, )
	assert.Nil(t, res.Data)
}

func TestSignInWithInvalidPassword(t *testing.T) {
	body, _ := json.Marshal(&auth.SignInParams{
		Account:  tester.Username,
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
	body, _ := json.Marshal(&auth.SignInParams{
		Account:  tester.Username,
		Password: tester.Password,
	})

	// empty body
	r := tester.Http.Post("/v1/auth/signin", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := auth.SignInResponse{}

	if err := mapstructure.Decode(res.Data, &profile); err != nil {
		assert.Error(t, err, err.Error())
	}

	assert.NotEmpty(t, profile.Token)

	if c, err := token.Parse(profile.Token); err != nil {
		t.Error(err)
		return
	} else {
		assert.IsType(t, c.Uid, int64(123), "UID必须是int64")
	}
}
