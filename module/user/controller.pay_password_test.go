// Copyright 2019 Axetroy. All rights reserved. MIT license.
package user_test

import (
	"github.com/axetroy/go-server/exception"
	"github.com/axetroy/go-server/module/auth"
	"github.com/axetroy/go-server/module/user"
	"github.com/axetroy/go-server/module/user/user_error"
	"github.com/axetroy/go-server/module/user/user_schema"
	"github.com/axetroy/go-server/schema"
	"github.com/axetroy/go-server/service/redis"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestSetPayPassword(t *testing.T) {
	var (
		testUser user_schema.Profile
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

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = user_schema.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer auth.DeleteUserByUserName(username)

		assert.False(t, testUser.PayPassword)
	}

	{
		// 2. 设置交易密码失败
		r := user.SetPayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123", // 两次密码不一致
			PasswordConfirm: "321321",
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, user_error.ErrInvalidConfirmPassword.Error(), r.Message)
	}

	{
		// 3. 设置交易密码成功
		r := user.SetPayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123",
			PasswordConfirm: "123123",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.True(t, r.Data.(bool))
	}

	{
		// 4. 已经设置过了，再设置就报错
		r := user.SetPayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123",
			PasswordConfirm: "123123",
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.ErrPayPasswordSet.Error(), r.Message)
		assert.False(t, r.Data.(bool))
	}
}

func TestUpdatePayPassword(t *testing.T) {
	var (
		testUser user_schema.Profile
	)

	{
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestUpdatePayPassword"
		password := "123123"

		r := auth.SignUp(auth.SignUpParams{
			Username: &username,
			Password: password,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = user_schema.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer auth.DeleteUserByUserName(username)

		assert.False(t, testUser.PayPassword)
	}

	{
		// 2. 更新交易密码失败, 因为此时还没有交易密码
		r := user.UpdatePayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.UpdatePayPasswordParams{
			OldPassword: "321321",
			NewPassword: "123123",
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, user_error.ErrErrRequireErrPayPasswordSet.Error(), r.Message)
		assert.False(t, r.Data.(bool))
	}

	{
		// 3. 设置交易密码成功
		r := user.SetPayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123",
			PasswordConfirm: "123123",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.True(t, r.Data.(bool))
	}

	{
		// 4. 更新交易密码失败, 原密码错误
		r := user.UpdatePayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.UpdatePayPasswordParams{
			OldPassword: "321321",
			NewPassword: "111111",
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, exception.ErrInvalidPassword.Error(), r.Message)
		assert.False(t, r.Data.(bool))
	}

	{
		// 5. 更新交易密码成功
		r := user.UpdatePayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.UpdatePayPasswordParams{
			OldPassword: "123123",
			NewPassword: "321321",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.True(t, r.Data.(bool))
	}
}

func TestResetPayPassword(t *testing.T) {
	var (
		testUser user_schema.Profile
	)

	{
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestResetPayPassword"
		password := "123123"

		r := auth.SignUp(auth.SignUpParams{
			Username: &username,
			Password: password,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = user_schema.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer auth.DeleteUserByUserName(username)

		assert.False(t, testUser.PayPassword)
	}

	// 生成重置码
	resetCode := user.GenerateResetPayPasswordCode(testUser.Id)

	// redis缓存重置码
	assert.Nil(t, redis.ResetCodeClient.Set(resetCode, testUser.Id, time.Minute*10).Err())

	{
		// 2. 重置交易密码失败, 因为此时还没有交易密码
		r := user.ResetPayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.ResetPayPasswordParams{
			Code:        resetCode,
			NewPassword: "123123",
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, user_error.ErrErrRequireErrPayPasswordSet.Error(), r.Message)
		assert.False(t, r.Data.(bool))
	}

	{
		// 3. 设置交易密码成功
		r := user.SetPayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123",
			PasswordConfirm: "123123",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.True(t, r.Data.(bool))
	}

	{
		// 4. 重置交易密码失败, 错误的重置码
		r := user.ResetPayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.ResetPayPasswordParams{
			Code:        "123123",
			NewPassword: "123123",
		})

		assert.Equal(t, schema.StatusFail, r.Status)
		assert.Equal(t, user_error.ErrInvalidResetCode.Error(), r.Message)
		assert.False(t, r.Data.(bool))
	}

	{
		// 5. 重置交易密码成功
		r := user.ResetPayPassword(schema.Context{
			Uid: testUser.Id,
		}, user.ResetPayPasswordParams{
			Code:        resetCode,
			NewPassword: "123123",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.True(t, r.Data.(bool))
	}
}
