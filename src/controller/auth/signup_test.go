// Copyright 2019 Axetroy. All rights reserved. MIT license.
package auth_test

import (
	"encoding/json"
	"github.com/axetroy/go-server/src/controller/auth"
	"github.com/axetroy/go-server/src/controller/invite"
	"github.com/axetroy/go-server/src/exception"
	"github.com/axetroy/go-server/src/model"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestSignUpWithEmptyBody(t *testing.T) {
	// empty body
	r := tester.HttpUser.Post("/v1/auth/signup", []byte(nil), nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes()), &res))

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, exception.InvalidParams.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestSignUpWithNotFullBody(t *testing.T) {
	username := "username"

	// 没有输入密码
	body, _ := json.Marshal(&auth.SignUpParams{
		Username: &username,
	})

	// empty body
	r := tester.HttpUser.Post("/v1/auth/signup", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := schema.Response{}

	assert.Nil(t, json.Unmarshal(r.Body.Bytes()), &res))

	assert.Equal(t, schema.StatusFail, res.Status)
	assert.Equal(t, exception.RequirePassword.Error(), res.Message)
	assert.Nil(t, res.Data)
}

func TestSignUpSuccess(t *testing.T) {
	rand.Seed(99) // 重置随机码，否则随机数会一样

	username := "test-TestSignUpSuccess"

	res := auth.SignUp(auth.SignUpParams{
		Username: &username,
		Password: "123123",
	})

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	defer auth.DeleteUserByUserName(username)

	profile := schema.Profile{}

	assert.Nil(t, tester.Decode(res.Data, &profile))

	// 默认未激活状态
	assert.Equal(t, int(profile.Status), int(model.UserStatusInactivated))
	assert.Equal(t, profile.Username, username)
	assert.Equal(t, *profile.Nickname, username)
	assert.Equal(t, profile.Role, []string{model.DefaultUser.Name})
	assert.Nil(t, profile.Email)
	assert.Nil(t, profile.Phone)
}

func TestSignUpInviteCode(t *testing.T) {
	rand.Seed(133) // 重置随机码，否则随机数会一样

	testerUsername := "tester"
	testerUid := ""
	username := "test-TestSignUpInviteCode"

	inviteCode := ""

	// 动态创建一个测试账号
	{
		r := auth.SignUp(auth.SignUpParams{
			Username: &testerUsername,
			Password: "123123",
		})

		profile := schema.Profile{}

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		assert.Nil(t, tester.Decode(r.Data, &profile))

		inviteCode = profile.InviteCode
		testerUid = profile.Id

		defer auth.DeleteUserByUserName(testerUsername)
	}

	rand.Seed(1111) // 重置随机码，否则随机数会一样

	res := auth.SignUp(auth.SignUpParams{
		Username:   &username,
		Password:   "123123",
		InviteCode: &inviteCode,
	})

	assert.Equal(t, schema.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	defer auth.DeleteUserByUserName(username)

	profile := schema.Profile{}

	if !assert.Nil(t, tester.Decode(res.Data, &profile)) {
		return
	}

	// 默认未激活状态
	assert.Equal(t, int(model.UserStatusInactivated), int(profile.Status))
	assert.Equal(t, username, profile.Username)
	assert.Equal(t, username, *profile.Nickname)
	assert.Nil(t, profile.Email)
	assert.Nil(t, profile.Phone)

	// 获取我的邀请记录
	resInvite := invite.GetByStruct(&model.InviteHistory{Invitee: profile.Id})
	InviteeData := schema.Invite{}

	assert.Nil(t, tester.Decode(resInvite.Data, &InviteeData))
	assert.Equal(t, profile.Id, InviteeData.Invitee)
	assert.Equal(t, testerUid, InviteeData.Inviter)
}
