package auth_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/model"
	"github.com/axetroy/go-server/orm"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/router"
	"github.com/axetroy/go-server/tester"
	"github.com/axetroy/mocker"
	"github.com/go-xorm/xorm"
	"github.com/stretchr/testify/assert"
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
	rand.Seed(99) // 重置随机码，否则随机数会一样
	m := mocker.New(router.Router)

	username := "test-TestSignUpSuccess"

	body, _ := json.Marshal(&auth.SignUpParams{
		Username: &username,
		Password: "123123",
	})

	// empty body
	r := m.Post("/v1/auth/signup", body, nil)

	var (
		session *xorm.Session
		ormErr  error
	)

	session = orm.Db.NewSession()

	if ormErr = session.Begin(); ormErr != nil {
		return
	}

	defer func() {
		if ormErr != nil {
			_ = session.Rollback()
		} else {
			_ = session.Commit()
		}
	}()

	assert.Equal(t, http.StatusOK, r.Code)

	res := response.Response{}

	assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res))

	assert.Equal(t, response.StatusSuccess, res.Status)
	assert.Equal(t, "", res.Message)

	profile := user.Profile{}

	if assert.Nil(t, tester.Decode(res.Data, &profile)) {
		return
	}

	defer func() {
		raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE id = %v", "user", profile.Id)

		if _, err := session.Exec(raw); err != nil {
			t.Error(err)
			return
		} else {

		}
	}()

	// 默认未激活状态
	assert.Equal(t, int(profile.Status), int(model.UserStatusInactivated))
	assert.Equal(t, profile.Username, username)
	assert.Equal(t, *profile.Nickname, username)
	assert.Nil(t, profile.Email)
	assert.Nil(t, profile.Phone)

}

func TestSignUpInviteCode(t *testing.T) {
	rand.Seed(100) // 重置随机码，否则随机数会一样
	m := mocker.New(router.Router)

	username := "test-TestSignUpInviteCode"

	inviteCode := "11111111"

	body, _ := json.Marshal(&auth.SignUpParams{
		Username:   &username,
		Password:   "123123",
		InviteCode: &inviteCode,
	})

	// empty body
	r := m.Post("/v1/auth/signup", body, nil)

	var (
		session *xorm.Session
		ormErr  error
	)

	session = orm.Db.NewSession()

	if ormErr = session.Begin(); ormErr != nil {
		return
	}

	defer func() {
		if ormErr != nil {
			fmt.Println(ormErr)
			_ = session.Rollback()
		} else {
			_ = session.Commit()
		}
	}()

	if !assert.Equal(t, http.StatusOK, r.Code) {
		return
	}

	res := response.Response{}

	if !assert.Nil(t, json.Unmarshal([]byte(r.Body.String()), &res)) {
		return
	}
	if !assert.Equal(t, res.Status, response.StatusSuccess) {
		fmt.Println(res.Message)
		return
	}
	if !assert.Equal(t, res.Message, "") {
		return
	}

	profile := user.Profile{}

	if !assert.Nil(t, tester.Decode(res.Data, &profile)) {
		return
	}

	defer func() {
		raw := fmt.Sprintf("DELETE FROM \"%v\" WHERE id = '%s'", "user", profile.Id)

		if _, err := session.Exec(raw); err != nil {
			ormErr = err
			return
		}

		raw1 := fmt.Sprintf("DELETE FROM \"%v\" WHERE invited = '%s'", "invite_history", profile.Id)

		if _, err := session.Exec(raw1); err != nil {
			ormErr = err
			return
		}
	}()

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

	var inviteHistory model.InviteHistory

	if isExist, er := session.Where("invited = ?", profile.Id).Get(&inviteHistory); er != nil {
		ormErr = er
	} else {
		if isExist == false {
			t.Error(errors.New("根据邀请码注册的，应该会产生一条邀请记录"))
		}
	}

	if !assert.Equal(t, profile.Id, inviteHistory.Invited) {
		return
	}
	if !assert.Equal(t, tester.Uid, inviteHistory.Invitor) {
		return
	}
}
