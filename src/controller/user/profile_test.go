package user_test

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/user"
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

	r := tester.HttpUser.Get("/v1/user/profile", []byte(""), &header)

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

	r := tester.HttpUser.Get("/v1/user/profile", []byte(""), &header)

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
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	r := user.GetProfile(controller.Context{Uid: userInfo.Id})

	profile := schema.Profile{}

	assert.Nil(t, tester.Decode(r.Data, &profile))

	assert.Equal(t, userInfo.Id, profile.Id)
	assert.Equal(t, userInfo.Username, profile.Username)
	assert.Equal(t, userInfo.CreatedAt, profile.CreatedAt)
}

func TestGetProfileByAdmin(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	r := user.GetProfileByAdmin(controller.Context{Uid: adminInfo.Id}, userInfo.Id)

	profile := schema.Profile{}

	assert.Nil(t, tester.Decode(r.Data, &profile))

	assert.Equal(t, userInfo.Id, profile.Id)
	assert.Equal(t, userInfo.Username, profile.Username)
	assert.Equal(t, userInfo.CreatedAt, profile.CreatedAt)
}

func TestGetProfileRouter(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + userInfo.Token,
	}

	r := tester.HttpUser.Get("/v1/user/profile", []byte(""), &header)

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

	if !assert.Equal(t, userInfo.Id, profile.Id) {
		return
	}
	if !assert.Equal(t, userInfo.Username, *profile.Email) {
		return
	}
}

func TestGetProfileByAdminRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": util.TokenPrefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/user/u/"+userInfo.Id, nil, &header)

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

	if !assert.Equal(t, userInfo.Id, profile.Id) {
		return
	}
	if !assert.Equal(t, userInfo.Username, *profile.Email) {
		return
	}
}
