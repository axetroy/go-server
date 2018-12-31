package user_test

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
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
		"Authorization": "Bearer 12312",
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
	header := mocker.Header{
		"Authorization": tester.Token,
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
		return
	}
	if !assert.Equal(t, "", res.Message) {
		return
	}

	profile := user.Profile{}

	if assert.Nil(t, tester.Decode(res.Data, &profile)) {
		return
	}

	fmt.Println(profile)

	if !assert.Equal(t, tester.Uid, profile.Id) {
		return
	}
	if !assert.Equal(t, tester.Username, *profile.Email) {
		return
	}
}
