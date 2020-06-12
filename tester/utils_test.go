package tester_test

import (
	"github.com/axetroy/go-server/internal/app/admin_server/controller/admin"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"testing"
)

func init() {
	admin.CreateAdmin(admin.CreateAdminParams{
		Account:  "admin",
		Password: "123456",
		Name:     "admin",
	}, true)
}

func TestCreateUser(t *testing.T) {
	user, err := tester.CreateUser()

	assert.Nil(t, err)

	defer tester.DeleteUserByUserName(user.Username)

	assert.NotEmpty(t, user.Username)
	assert.NotEmpty(t, user.Token)
	assert.NotEmpty(t, user.Id)
}

func TestLoginAdmin(t *testing.T) {
	adminInfo, err := tester.LoginAdmin()

	assert.Nil(t, err)

	assert.NotEmpty(t, adminInfo.Token)
	assert.NotEmpty(t, adminInfo.Id)
	assert.Equal(t, "admin", adminInfo.Username)
	assert.Equal(t, "admin", adminInfo.Name)
}

func TestCreateWaiter(t *testing.T) {
	waiterInfo, err := tester.CreateWaiter()

	assert.Nil(t, err)

	assert.NotEmpty(t, waiterInfo.Token)
	assert.NotEmpty(t, waiterInfo.Id)
}
