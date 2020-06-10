// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package user_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/internal/app/admin_server/controller/user"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/token"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestUpdateProfileByAdmin(t *testing.T) {
	var (
		nickName = "nickname"
	)
	userInfo, _ := tester.CreateUser()
	adminInfo, _ := tester.LoginAdmin()

	defer tester.DeleteUserByUserName(userInfo.Username)

	res := user.UpdateProfileByAdmin(helper.Context{Uid: adminInfo.Id}, userInfo.Id, user.UpdateProfileParams{
		Nickname: &nickName,
	})

	profile := schema.Profile{}

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)
	assert.Nil(t, res.Decode(&profile))

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

	defer tester.DeleteUserByUserName(userInfo.Username)

	header := mocker.Header{
		"Authorization": token.Prefix + " " + adminInfo.Token,
	}

	body, _ := json.Marshal(&user.UpdateProfileParams{
		Nickname: &nickName,
	})

	r := tester.HttpAdmin.Put("/v1/user/"+userInfo.Id, body, &header)

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
	assert.Nil(t, res.Decode(&profile))

	assert.Equal(t, nickName, *profile.Nickname)
}
