package login_test

import (
	loginLog "github.com/axetroy/go-server/internal/app/admin_server/controller/logger/login"
	"github.com/axetroy/go-server/internal/library/helper"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetLatestLoginLog(t *testing.T) {
	user, _ := tester.CreateUser()

	defer tester.DeleteUserByUserName(user.Username)

	r := loginLog.GetLatestLoginLog(helper.Context{
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
