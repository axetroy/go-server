// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package ws_test

import (
	"fmt"
	"github.com/axetroy/go-server/internal/app/customer_service/ws"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMatcher_Join(t *testing.T) {
	// 没有客服的时候
	// 应该在排队
	matcher := ws.NewMatcher()

	matcher.SetWaiterClientFunc(func(id string) *ws.Client {
		return &ws.Client{
			UUID:  "waiter",
			Ready: true,
		}
	})

	{
		waiter, location := matcher.Join("test")

		assert.Nil(t, waiter)
		assert.Equal(t, 0, location)

		assert.Equal(t, 1, matcher.GetPendingLength())

		assert.Equal(t, "test", *matcher.ShiftPending())
		assert.Equal(t, 0, matcher.GetPendingLength())
	}

	// 测试离开
	{
		_, _ = matcher.Join("test")

		assert.Equal(t, 1, matcher.GetPendingLength())

		matcher.Leave("test")

		assert.Equal(t, 0, matcher.GetPendingLength())
	}

	// 如果客服已存在
	{
		matcher.AddWaiter("waiter")
		matcher.Join("user1")

		assert.Equal(t, map[string][]string{
			"waiter": {"user1"},
		}, matcher.GetMatcher())

		matcher.Join("user2")
		assert.Equal(t, map[string][]string{
			"waiter": {"user1", "user2"},
		}, matcher.GetMatcher())

		matcher.Leave("user1")
		matcher.Leave("user2")
	}

	// 测试
	{
		matcher = ws.NewMatcher()
		matcher.SetWaiterClientFunc(func(id string) *ws.Client {
			return &ws.Client{
				UUID:  "waiter",
				Ready: true,
			}
		})
		matcher.AddWaiter("waiter")

		index := 0

		for {
			if index > matcher.Max+1 {
				break
			}

			_, _ = matcher.Join(fmt.Sprintf("%d", index))

			index = index + 1
		}

		assert.Equal(t, map[string][]string{
			"waiter": {"0", "1", "2", "3", "4"},
		}, matcher.GetMatcher())

		assert.Equal(t, []string{"5", "6"}, matcher.GetPendingQueue())

	}
}

func TestMatcher_AddWaiter(t *testing.T) {
	matcher := ws.NewMatcher()

	// 添加客服
	{
		matcher.AddWaiter("test")

		assert.Equal(t, map[string][]string{
			"test": {},
		}, matcher.GetMatcher())
	}

	// 再添加相同的客服
	{
		matcher.AddWaiter("test")

		assert.Equal(t, map[string][]string{
			"test": {},
		}, matcher.GetMatcher())
	}

	// 客服离开
	{
		matcher.RemoveWaiter("test")

		assert.Equal(t, map[string][]string{}, ws.MatcherPool.GetMatcher())
	}
}

// 客服移除，那么剩下的用户会分配到队列中
func TestMatcher_RemoveWaiter(t *testing.T) {
	matcher := ws.NewMatcher()

	matcher.SetWaiterClientFunc(func(id string) *ws.Client {
		return &ws.Client{
			UUID:  "test",
			Ready: true,
		}
	})

	// 添加客服
	matcher.AddWaiter("test")

	matcher.Join("user1")
	matcher.Join("user2")

	assert.Equal(t, []string{}, matcher.GetPendingQueue())

	matcher.RemoveWaiter("test")

	assert.Equal(t, []string{"user2", "user1"}, matcher.GetPendingQueue())
}
