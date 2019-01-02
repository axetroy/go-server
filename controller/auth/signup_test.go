package auth_test

import (
	"encoding/json"
	"fmt"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/controller/invite"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"net/http"
	"testing"
)

func TestSignUpWithEmptyBody(t *testing.T) {
	// empty body
	r := tester.Http.Post("/v1/auth/signup", []byte(nil), nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, res.Status, response.StatusFail)
	assert.Equal(t, res.Message, exception.InvalidParams.Error())
	assert.Nil(t, res.Data)
}

func TestSignUpWithNotFullBody(t *testing.T) {
	username := "username"

	// 没有输入密码
	body, _ := json.Marshal(&auth.SignUpParams{
		Username: &username,
	})

	// empty body
	r := tester.Http.Post("/v1/auth/signup", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, res.Status, response.StatusFail)
	assert.Equal(t, res.Message, exception.RequirePassword.Error())
	assert.Nil(t, res.Data)
}

func TestSignUpSuccess(t *testing.T) {
	rand.Seed(99) // 重置随机码，否则随机数会一样

	username := "test-TestSignUpSuccess"

	res := auth.SignUp(auth.SignUpParams{
		Username: &username,
		Password: "123123",
	})

	if !assert.Equal(t, res.Status, response.StatusSuccess) {
		fmt.Println(res.Message)
		return
	}

	if !assert.Equal(t, res.Message, "") {
		return
	}

	defer func() {
		auth.DeleteUserByUserName(username)
	}()

	profile := user.Profile{}

	if assert.Nil(t, tester.Decode(res.Data, &profile)) {
		return
	}

	// 默认未激活状态
	assert.Equal(t, int(profile.Status), int(model.UserStatusInactivated))
	assert.Equal(t, profile.Username, username)
	assert.Equal(t, *profile.Nickname, username)
	assert.Nil(t, profile.Email)
	assert.Nil(t, profile.Phone)

}

func TestSignUpInviteCode(t *testing.T) {
	rand.Seed(133) // 重置随机码，否则随机数会一样

	username := "test-TestSignUpInviteCode"

	inviteCode := tester.InviteCode

	res := auth.SignUp(auth.SignUpParams{
		Username: &username,
		Password: "123123",
		// TODO: 动态生成一个账号
		InviteCode: &inviteCode,
	})

	if !assert.Equal(t, res.Status, response.StatusSuccess) {
		fmt.Println(res.Message)
		return
	}
	if !assert.Equal(t, res.Message, "") {
		return
	}

	defer func() {
		defer func() {
			auth.DeleteUserByUserName(username)
		}()
	}()

	profile := user.Profile{}

	if !assert.Nil(t, tester.Decode(res.Data, &profile)) {
		return
	}

	// 默认未激活状态
	assert.Equal(t, int(profile.Status), int(model.UserStatusInactivated))
	assert.Equal(t, profile.Username, username)
	assert.Equal(t, *profile.Nickname, username)
	if !assert.Nil(t, profile.Email) {
		return
	}
	if !assert.Nil(t, profile.Phone) {
		return
	}

	// 获取我的邀请记录
	resInvite := invite.GetInviteById(&model.InviteHistory{Invited: profile.Id,})
	inviteData := invite.Invite{}

	if !assert.Nil(t, tester.Decode(resInvite.Data, &inviteData)) {
		return
	}

	if !assert.Equal(t, profile.Id, inviteData.Invited) {
		return
	}
	if !assert.Equal(t, "86303081515450368", inviteData.Invitor) {
		return
	}
}
