package login_test

import (
	"github.com/axetroy/go-server/src/controller"
	"github.com/axetroy/go-server/src/controller/auth"
	loginLog "github.com/axetroy/go-server/src/controller/logs/login"
	"github.com/axetroy/go-server/src/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLatestLoginLog(t *testing.T) {
	user, _ := tester.CreateUser()

	defer auth.DeleteUserByUserName(user.Username)

	r := loginLog.GetLatestLoginLog(controller.Context{
		Uid: user.Id,
	})

	assert.Equal(t, schema.StatusSuccess, r.Status)
	assert.Equal(t, "", r.Message)

	loginLogInfo := schema.LogLogin{}

	assert.Nil(t, tester.Decode(r.Data, &loginLogInfo))

	assert.Equal(t, user.Id, loginLogInfo.Uid)
	assert.Equal(t, user.Username, loginLogInfo.User.Username)
	assert.Equal(t, user.Nickname, loginLogInfo.User.Nickname)
	assert.Equal(t, user.Avatar, loginLogInfo.User.Avatar)
}
