// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package user_test

import (
	"github.com/axetroy/go-server/internal/controller"
	"github.com/axetroy/go-server/internal/controller/auth"
	"github.com/axetroy/go-server/internal/controller/user"
	"github.com/axetroy/go-server/internal/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/axetroy/go-server/internal/service/redis"
	"github.com/axetroy/go-server/tester"
	"github.com/stretchr/testify/assert"
	"math/rand"
	"testing"
	"time"
)

func TestSetPayPassword(t *testing.T) {
	var (
		testUser schema.Profile
	)

	{
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestSetPayPassword"
		password := "123123"

		r := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
			Username: username,
			Password: password,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = schema.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer auth.DeleteUserByUserName(username)

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

		assert.Equal(t, exception.InvalidConfirmPassword.Code(), r.Status)
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

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.Equal(t, nil, r.Data)
	}

	{
		// 4. 已经设置过了，再设置就报错
		r := user.SetPayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123",
			PasswordConfirm: "123123",
		})

		assert.Equal(t, exception.PayPasswordSet.Code(), r.Status)
		assert.Equal(t, exception.PayPasswordSet.Error(), r.Message)
		assert.Equal(t, nil, r.Data)
	}
}

func TestUpdatePayPassword(t *testing.T) {
	var (
		testUser schema.Profile
	)

	{
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestUpdatePayPassword"
		password := "123123"

		r := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
			Username: username,
			Password: password,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = schema.Profile{}

		if err := tester.Decode(r.Data, &testUser); err != nil {
			t.Error(err)
			return
		}

		defer auth.DeleteUserByUserName(username)

		assert.False(t, testUser.PayPassword)
	}

	{
		// 2. 更新交易密码失败, 因为此时还没有交易密码
		r := user.UpdatePayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.UpdatePayPasswordParams{
			OldPassword: "321321",
			NewPassword: "123123",
		})

		assert.Equal(t, exception.RequirePayPasswordSet.Code(), r.Status)
		assert.Equal(t, exception.RequirePayPasswordSet.Error(), r.Message)
		assert.Equal(t, nil, r.Data)
	}

	{
		// 3. 设置交易密码成功
		r := user.SetPayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123",
			PasswordConfirm: "123123",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.Equal(t, nil, r.Data)
	}

	{
		// 4. 更新交易密码失败, 原密码错误
		r := user.UpdatePayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.UpdatePayPasswordParams{
			OldPassword: "321321",
			NewPassword: "111111",
		})

		assert.Equal(t, exception.InvalidPassword.Code(), r.Status)
		assert.Equal(t, exception.InvalidPassword.Error(), r.Message)
		assert.Equal(t, nil, r.Data)
	}

	{
		// 5. 更新交易密码成功
		r := user.UpdatePayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.UpdatePayPasswordParams{
			OldPassword: "123123",
			NewPassword: "321321",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.Equal(t, nil, r.Data)
	}
}

func TestResetPayPassword(t *testing.T) {
	var (
		testUser schema.Profile
	)

	{
		// 1。 创建测试账号
		rand.Seed(111)
		username := "test-TestResetPayPassword"
		password := "123123"

		r := auth.SignUpWithUsername(auth.SignUpWithUsernameParams{
			Username: username,
			Password: password,
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)

		testUser = schema.Profile{}

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
	assert.Nil(t, redis.ClientResetCode.Set(resetCode, testUser.Id, time.Minute*10).Err())

	{
		// 2. 重置交易密码失败, 因为此时还没有交易密码
		r := user.ResetPayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.ResetPayPasswordParams{
			Code:        resetCode,
			NewPassword: "123123",
		})

		assert.Equal(t, exception.RequirePayPasswordSet.Code(), r.Status)
		assert.Equal(t, exception.RequirePayPasswordSet.Error(), r.Message)
		assert.Equal(t, nil, r.Data)
	}

	{
		// 3. 设置交易密码成功
		r := user.SetPayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.SetPayPasswordParams{
			Password:        "123123",
			PasswordConfirm: "123123",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.Equal(t, nil, r.Data)
	}

	{
		// 4. 重置交易密码失败, 错误的重置码
		r := user.ResetPayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.ResetPayPasswordParams{
			Code:        "123123",
			NewPassword: "123123",
		})

		assert.Equal(t, exception.InvalidResetCode.Code(), r.Status)
		assert.Equal(t, exception.InvalidResetCode.Error(), r.Message)
		assert.Equal(t, nil, r.Data)
	}

	{
		// 5. 重置交易密码成功
		r := user.ResetPayPassword(controller.Context{
			Uid: testUser.Id,
		}, user.ResetPayPasswordParams{
			Code:        resetCode,
			NewPassword: "123123",
		})

		assert.Equal(t, schema.StatusSuccess, r.Status)
		assert.Equal(t, "", r.Message)
		assert.Equal(t, nil, r.Data)
	}
}
