// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package ws

import (
	"sync"
)

type Matcher struct {
	sync.RWMutex
	Broadcast chan bool           // 调度器，当收到通知时，就安排客服接待排队的用户
	Max       int                 // 一个客服最多接待多少个用户
	matcher   map[string][]string // 已经匹配的 socket对
	pending   []string            // 排队的用户 socket
}

func NewMatcher() *Matcher {
	return &Matcher{
		Max:       5, // 一个客服最多接待 5 个用户
		matcher:   map[string][]string{},
		Broadcast: make(chan bool),
	}
}

var MatcherPool = NewMatcher()

func (c *Matcher) GetPendingLength() int {
	c.RLock()
	defer c.RUnlock()
	return len(c.pending)
}

func (c *Matcher) GetPendingQueue() []string {
	return c.pending
}

func (c *Matcher) ShiftPending() *string {
	c.RLock()
	defer c.RUnlock()
	if len(c.pending) == 0 {
		return nil
	}
	userSocketUUID := c.pending[len(c.pending)-1]

	c.pending = c.pending[1:]

	return &userSocketUUID
}

func (c *Matcher) GetMatcher() map[string][]string {
	return c.matcher
}

// 用户加入匹配池
// 返回接待的客服 UUID
// 如果返回空，那么说明没有找到合适的客服，加入等待队列
// 第二个参数用于插入到最前面的队列，出于最优先级
// 返回 int 代表出于队列的地 n 位
func (c *Matcher) Join(userSocketUUID string, prepend ...bool) (*string, int) {
	c.RLock()
	defer c.RUnlock()
	idleWaiter := c.GetIdleWaiter()

	// 如果找不到最佳的客服，那么先加入队列
	if idleWaiter == nil {
		// 确保当前连接不在队列中
		for index, id := range c.pending {
			if id == userSocketUUID {
				return nil, index
			}
		}

		if len(prepend) > 0 && prepend[0] {
			c.pending = append([]string{userSocketUUID}, c.pending...)
		} else {
			c.pending = append(c.pending, userSocketUUID)
		}
		return nil, len(c.pending) - 1
	} else {
		// 如果找到了空闲的客服
		for waiter, users := range c.matcher {
			if waiter == *idleWaiter {
				c.matcher[waiter] = append(users, userSocketUUID)
				return idleWaiter, 0
			}
		}
	}

	return nil, 0
}

// 用户离开匹配池
func (c *Matcher) Leave(userSocketUUID string) {
	c.RLock()
	defer c.RUnlock()
	for waiter, users := range c.matcher {
		for index, user := range users {
			if user == userSocketUUID {
				c.matcher[waiter] = append(c.matcher[waiter][:index], c.matcher[waiter][index+1:]...)
			}
		}
	}

	for index, id := range c.pending {
		if id == userSocketUUID {
			c.pending = append(c.pending[:index], c.pending[index+1:]...)
		}
	}
}

// 获取这个客服当前服务的用户
func (c *Matcher) GetMyUsers(waiterSocketUUID string) []string {
	c.RLock()
	defer c.RUnlock()

	for id, users := range c.matcher {
		if id == waiterSocketUUID {
			if len(users) > c.Max {
				return users[:c.Max]
			} else {
				return users
			}
		}
	}

	return []string{}
}

// 添加客服
func (c *Matcher) AddWaiter(waiterSocketUUID string) {
	c.RLock()
	defer c.RUnlock()

	for id := range c.matcher {
		if id == waiterSocketUUID {
			return
		}
	}

	c.matcher[waiterSocketUUID] = []string{}

	// 如果这时候等待队列里面有排队的，就先处理它
	if len(c.pending) > 0 {
		var users []string
		if len(c.pending) > c.Max {
			users = c.pending[:c.Max]
			c.pending = c.pending[c.Max:]
		} else {
			users = c.pending
			c.pending = []string{}
		}

		for _, userSocketUUID := range users {
			c.Join(userSocketUUID)
		}
	}
}

// 移除客服
func (c *Matcher) RemoveWaiter(waiterSocketUUID string) {
	c.RLock()
	defer c.RUnlock()

	for id, users := range c.matcher {
		if id == waiterSocketUUID {
			delete(c.matcher, id)

			// 还出于连接的用户，放入到队列中
			// 并且优先放在第一排
			for _, user := range users {
				c.Join(user, true)
			}
		}
	}
}

// 获取当前最空闲的客服
func (c *Matcher) GetIdleWaiter() *string {
	c.RLock()
	defer c.RUnlock()

	var (
		bestWaiterId      *string
		currentUserNumber = c.Max
	)

	for waiter, users := range c.matcher {
		if len(users) == 0 {
			bestWaiterId = &waiter
			break
		} else if len(users) < currentUserNumber {
			bestWaiterId = &waiter
		}
	}

	return bestWaiterId
}

// 获取当前接待我的客服
func (c *Matcher) GetMyWaiter(userSocketUUID string) *string {
	c.RLock()
	defer c.RUnlock()

	for id, users := range c.matcher {
		for _, u := range users {
			if u == userSocketUUID {
				return &id
			}
		}
	}

	return nil
}
