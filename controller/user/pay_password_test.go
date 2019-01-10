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

func TestSetPayPassword(t *testing.T) {
	var (
		testUser user.Profile
	)

	{
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestSetPayPassword"
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

		assert.False(t, testUser.PayPassword)
	}

	{
		// 2. 设置交易密码失败
		r := user.SetPayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123", // 两次密码不一致
			PasswordConfirm: "321321",
		})

		assert.Equal(t, response.StatusFail, r.Status)
		assert.Equal(t, exception.InvalidConfirmPassword.Error(), r.Message)
	}

	{
		// 3. 设置交易密码成功
		r := user.SetPayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123",
			PasswordConfirm: "123123",
		})

		assert.Equal(t, response.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.True(t, r.Data.(bool))
	}

	{
		// 4. 已经设置过了，再设置就报错
		r := user.SetPayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123",
			PasswordConfirm: "123123",
		})

		assert.Equal(t, response.StatusFail, r.Status)
		assert.Equal(t, exception.PayPasswordSet.Error(), r.Message)
		assert.False(t, r.Data.(bool))
	}
}
