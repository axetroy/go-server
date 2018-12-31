package auth_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axetroy/mocker"
	"github.com/go-xorm/xorm"
	"github.com/mitchellh/mapstructure"
	"github.com/stretchr/testify/assert"
	"gitlab.com/axetroy/server/controller/auth"
	"gitlab.com/axetroy/server/controller/user"
	"gitlab.com/axetroy/server/exception"
	"gitlab.com/axetroy/server/model"
	"gitlab.com/axetroy/server/orm"
	"gitlab.com/axetroy/server/response"
	"gitlab.com/axetroy/server/router"
	"gitlab.com/axetroy/server/tester"
	"math/rand"
	"net/http"
	"testing"
)

func RandomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) //A=65 and Z = 65+25
	}
	return string(bytes)
}

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
	m := mocker.New(router.Router)

	username := "test-" + RandomString(10)

	body, _ := json.Marshal(&auth.SignUpParams{
		Username: &username,
		Password: "123123",
	})

	// empty body
	r := m.Post("/v1/auth/signup", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, res.Status, response.StatusSuccess)
	assert.Equal(t, res.Message, "")

	profile := user.Profile{}

	if err := mapstructure.Decode(res.Data, &profile); err != nil {
		assert.Error(t, err, err.Error())
	}

	// 默认未激活状态
	assert.Equal(t, int(profile.Status), int(model.UserStatusInactivated))
	assert.Equal(t, profile.Username, username)
	assert.Equal(t, *profile.Nickname, username)
	assert.Nil(t, profile.Email)
	assert.Nil(t, profile.Phone)

	var (
		session *xorm.Session
		err     error
	)

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = session.Rollback()
		} else {
			_ = session.Commit()
		}
	}()

	raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE id = %v", "user", profile.Id)

	if _, err := session.Exec(raw); err != nil {
		t.Error(err)
		return
	} else {

	}
}

func TestSignUpInviteCode(t *testing.T) {
	m := mocker.New(router.Router)

	username := "test-" + RandomString(10)

	inviteCode := "11111111"

	body, _ := json.Marshal(&auth.SignUpParams{
		Username:   &username,
		Password:   "123123",
		InviteCode: &inviteCode,
	})

	// empty body
	r := m.Post("/v1/auth/signup", body, nil)

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, res.Status, response.StatusSuccess)
	assert.Equal(t, res.Message, "")

	profile := user.Profile{}

	if err := mapstructure.Decode(res.Data, &profile); err != nil {
		assert.Error(t, err, err.Error())
	}

	// 默认未激活状态
	assert.Equal(t, int(profile.Status), int(model.UserStatusInactivated))
	assert.Equal(t, profile.Username, username)
	assert.Equal(t, *profile.Nickname, username)
	assert.Nil(t, profile.Email)
	assert.Nil(t, profile.Phone)

	var (
		session *xorm.Session
		err     error
	)

	session = orm.Db.NewSession()

	if err = session.Begin(); err != nil {
		return
	}

	defer func() {
		if err != nil {
			_ = session.Rollback()
		} else {
			_ = session.Commit()
		}
	}()

	inviteHistory := model.InviteHistory{
		Invited: profile.Id,
	}

	if isExist, er := session.Get(&inviteHistory); er != nil {
		err = er
		t.Error(err)
	} else {
		if isExist == false {
			t.Error(errors.New("根据邀请码注册的，应该会产生一条邀请记录"))
		}
	}

	assert.Equal(t, profile.Id, inviteHistory.Invited)
	assert.Equal(t, tester.Uid, inviteHistory.Invitor)

	raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE id = %v", "user", profile.Id)

	if _, err := session.Exec(raw); err != nil {
		t.Error(err)
		return
	}

	raw1 := fmt.Sprintf("DELETE FROM \"%v\" WHERE invited = %v", "invite_history", profile.Id)

	if _, err := session.Exec(raw1); err != nil {
		t.Error(err)
		return
	}
}
