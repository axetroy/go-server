// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/core/controller"
	"github.com/axetroy/go-server/core/controller/auth"
	"github.com/axetroy/go-server/core/controller/user"
	"github.com/axetroy/go-server/core/exception"
	"github.com/axetroy/go-server/core/schema"
	"github.com/axetroy/go-server/core/service/token"
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

	r := tester.HttpUser.Get("/v1/user/profile", nil, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, exception.InvalidToken.Code(), res.Status)
	assert.Equal(t, exception.InvalidAuth.Error(), res.Message)
}

func TestGetProfileWithInvalidToken(t *testing.T) {
	header := mocker.Header{
		"Authorization": token.Prefix + " 12312",
	}

	r := tester.HttpUser.Get("/v1/user/profile", nil, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, exception.InvalidToken.Code(), res.Status)
	assert.Equal(t, exception.InvalidToken.Error(), res.Message)
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
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	r := tester.HttpUser.Get("/v1/user/profile", nil, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := schema.Profile{}

	assert.Nil(t, tester.Decode(res.Data, &profile))
	assert.Equal(t, userInfo.Id, profile.Id)
	assert.Equal(t, userInfo.Username, profile.Username)
}

func TestGetProfileByAdminRouter(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	r := tester.HttpAdmin.Get("/v1/user/u/"+userInfo.Id, nil, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := schema.Profile{}

	assert.Nil(t, tester.Decode(res.Data, &profile))
	assert.Equal(t, userInfo.Id, profile.Id)
	assert.Equal(t, userInfo.Username, profile.Username)
}

func TestUpdateProfile(t *testing.T) {
	var (
		nickName = "nickname"
	)
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	r := user.UpdateProfile(controller.Context{Uid: userInfo.Id}, user.UpdateProfileParams{
		Nickname: &nickName,
	})

	profile := schema.Profile{}

	assert.Nil(t, tester.Decode(r.Data, &profile))

	assert.Equal(t, userInfo.Id, profile.Id)
	assert.Equal(t, userInfo.Username, profile.Username)
	assert.Equal(t, nickName, *profile.Nickname)
	assert.Equal(t, userInfo.CreatedAt, profile.CreatedAt)
}

func TestUpdateProfileRouter(t *testing.T) {
	var (
		nickName = "nickname"
	)
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + userInfo.Token,
	}

	body, _ := json.Marshal(&user.UpdateProfileParams{
		Nickname: &nickName,
	})

	r := tester.HttpUser.Put("/v1/user/profile", body, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := schema.Profile{}

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)
	assert.Nil(t, tester.Decode(res.Data, &profile))

	assert.Equal(t, nickName, *profile.Nickname)
}

func TestUpdateProfileByAdmin(t *testing.T) {
	var (
		nickName = "nickname"
	)
	userInfo, _ := tester.CreateUser()
	adminInfo, _ := tester.LoginAdmin()

	defer auth.DeleteUserByUserName(userInfo.Username)

	res := user.UpdateProfileByAdmin(controller.Context{Uid: adminInfo.Id}, userInfo.Id, user.UpdateProfileParams{
		Nickname: &nickName,
	})

	profile := schema.Profile{}

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)
	assert.Nil(t, tester.Decode(res.Data, &profile))

	assert.Equal(t, userInfo.Id, profile.Id)
	assert.Equal(t, userInfo.Username, profile.Username)
	assert.Equal(t, nickName, *profile.Nickname)
	assert.Equal(t, userInfo.CreatedAt, profile.CreatedAt)
}

func TestUpdateProfileByAdminRouter(t *testing.T) {
	var (
		nickName = "nickname"
	)
	userInfo, _ := tester.CreateUser()
	adminInfo, _ := tester.LoginAdmin()

	defer auth.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	body, _ := json.Marshal(&user.UpdateProfileParams{
		Nickname: &nickName,
	})

	r := tester.HttpAdmin.Put("/v1/user/u/"+userInfo.Id, body, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes(), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := schema.Profile{}

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)
	assert.Nil(t, tester.Decode(res.Data, &profile))

	assert.Equal(t, nickName, *profile.Nickname)
}
