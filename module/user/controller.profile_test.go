// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/user"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestGetProfileWithErrInvalidAuth(t *testing.T) {

	header := mocker.Header{
		"Authorization": "Bearera 12312", // invalid Bearera
	}

	r := tester.HttpUser.Get("/v1/user/profile", nil, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, token.ErrInvalidAuth.Error(), res.Message)
}

func TestGetProfileWithErrInvalidToken(t *testing.T) {
	header := mocker.Header{
		"Authorization": token.Prefix + " 12312",
	}

	r := tester.HttpUser.Get("/v1/user/profile", nil, &header)

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, token.ErrInvalidToken.Error(), res.Message)
}

func TestGetProfile(t *testing.T) {
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	r := user.GetProfile(schema.Context{Uid: userInfo.Id})

	profile := user_schema.Profile{}

	assert.Nil(t, tester.Decode(r.Data, &profile))

	assert.Equal(t, userInfo.Id, profile.Id)
	assert.Equal(t, userInfo.Username, profile.Username)
	assert.Equal(t, userInfo.CreatedAt, profile.CreatedAt)
}

func TestGetProfileByAdmin(t *testing.T) {
	adminInfo, _ := tester.LoginAdmin()
	userInfo, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(userInfo.Username)

	r := user.GetProfileByAdmin(schema.Context{Uid: adminInfo.Id}, userInfo.Id)

	profile := user_schema.Profile{}

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

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := user_schema.Profile{}

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

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := user_schema.Profile{}

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

	r := user.UpdateProfile(schema.Context{Uid: userInfo.Id}, user.UpdateProfileParams{
		Nickname: &nickName,
	})

	profile := user_schema.Profile{}

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

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := user_schema.Profile{}

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

	res := user.UpdateProfileByAdmin(schema.Context{Uid: adminInfo.Id}, userInfo.Id, user.UpdateProfileParams{
		Nickname: &nickName,
	})

	profile := user_schema.Profile{}

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

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))
	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := user_schema.Profile{}

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)
	assert.Nil(t, tester.Decode(res.Data, &profile))

	assert.Equal(t, nickName, *profile.Nickname)
}
