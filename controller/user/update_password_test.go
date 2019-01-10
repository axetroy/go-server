package user_test

import (
	"github.com/axetroy/go-server/controller"
	"github.com/axetroy/go-server/controller/auth"
	"github.com/axetroy/go-server/controller/user"
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/response"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
)

func TestUpdatePassword(t *testing.T) {
	var (
		testUser user.Profile
	)

	{
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestUpdatePassword"
		password := "123123"

		r := auth.SignUp(auth.SignUpParams{
			Username: &username,
			Password: password,
		})

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = user.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer func() {
			auth.DeleteUserByUserName(username)
		}()
	}

	{
		// 2. 更新测试账号的密码, 旧密码错误
		r := user.UpdatePassword(testUser.Id, user.UpdatePasswordParams{
			OldPassword: "321321",
			NewPassword: "aaa",
		})

		assert.Equal(t, response.StatusFail, r.Status)
		assert.Equal(t, exception.InvalidPassword.Error(), r.Message)
	}

	{
		var newPassword = "321321"
		// 2. 更新测试账号的密码
		r := user.UpdatePassword(testUser.Id, user.UpdatePasswordParams{
			OldPassword: "123123",
			NewPassword: newPassword,
		})

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.True(t, r.Data.(bool))

		r2 := auth.SignIn(controller.Context{
			UserAgent: "test",
			Ip:        "0.0.0.0.0",
		}, auth.SignInParams{
			Account:  testUser.Username,
			Password: newPassword,
		})

		assert.Equal(t, response.StatusSuccess, r2.Status)
		assert.Equal(t, "", r2.Message)
	}
}
