package auth_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/controller/email"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/services/redis"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestActivationWithEmptyBody(t *testing.T) {

	// empty body
	r := tester.Http.Post("/v1/auth/activation", []byte(nil), nil)

	if ok := assert.Equal(t, http.StatusOK, r.Code); !ok {
		return
	}

	res := response.Response{}

	if ok := assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)); !ok {
		return
	}

	assert.Equal(t, response.StatusFail, res.Status)
	assert.Equal(t, exception.InvalidParams.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestActivationWithInvalidCode(t *testing.T) {
	body, _ := json.Marshal(&auth.ActivationParams{
		Code: "123",
	})

	// empty body
	r := tester.Http.Post("/v1/auth/activation", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, res.Status, response.StatusFail)
	assert.Equal(t, res.Message, exception.InvalidActiveCode.Error())
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
		})

		profile := user.Profile{}

		assert.Nil(t, tester.Decode(r.Data, &profile))

		testerUid = profile.Id

		defer func() {
			auth.DeleteUserByUserName(testerUsername)
		}()
	}

	// generate activation code
	activationCode := email.GenerateActivationCode(testerUid)

	// set activationCode to redis
	if err := redis.ActivationCode.Set(activationCode, testerUid, time.Minute*30).Err(); err != nil {
		t.Error(err)
		return
	}

	defer func() {
		// remove activation code
		_ = redis.ActivationCode.Del(activationCode).Err()
	}()

	body, _ := json.Marshal(&auth.ActivationParams{
		Code: activationCode,
	})

	// empty body
	r := tester.Http.Post("/v1/auth/activation", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)
}
