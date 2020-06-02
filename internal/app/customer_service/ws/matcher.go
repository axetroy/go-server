package ws

import (
	"fmt"
	"sync"
)

type Matcher struct {
	sync.RWMutex
	max     int                 // 一个客服最多接待多少个用户
	matcher map[string][]string // 已经匹配的 socket对
	pending []string            // 排队的用户 socket
}

func NewMatcher() *Matcher {
	return &Matcher{
		max:     5, // 一个客服最多接待 5 个用户
		matcher: map[string][]string{},
	}
}

var MatcherPool = NewMatcher()

// 添加用户到等待队列
func (c *Matcher) AppendToQueue(userSocketUUID string) {
	c.RLock()
	defer c.RUnlock()

	c.pending = append(c.pending, userSocketUUID)

	return
}

func (c *Matcher) RemoveFromQueue(userSocketUUID string) {
	c.RLock()
	defer c.RUnlock()

	for index, id := range c.pending {
		if id == userSocketUUID {
			c.pending = append(c.pending[:index], c.pending[index+1:]...)
		}
	}

	return
}

// 添加客服
func (c *Matcher) GetUsersForWaiter(waiterSocketUUID string) []string {
	c.RLock()
	defer c.RUnlock()

	for id, users := range c.matcher {
		fmt.Println(users)
		if id == waiterSocketUUID {
			if len(users) > c.max {
				return users[:c.max]
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

	for id, _ := range c.matcher {
		if id == waiterSocketUUID {
			return
		}
	}

	c.matcher[waiterSocketUUID] = []string{}

	fmt.Println("发现排队的", c.pending)

	// 如果这时候等待队列里面有排队的，就先处理它
	if len(c.pending) > 0 {
		var users []string
		if len(c.pending) > c.max {
			users = c.pending[:c.max]
		} else {
			users = c.pending
		}

		for _, userSockerUUID := range users {
			c.Connect(waiterSocketUUID, userSockerUUID)
		}
	}

	return
}

func (c *Matcher) RemoveWaiter(waiterSocketUUID string) {
	c.RLock()
	defer c.RUnlock()

	for id, users := range c.matcher {
		if id == waiterSocketUUID {
			delete(c.matcher, id)

			// 还出于连接的用户，放入到队列中
			// 并且优先放在第一排
			c.pending = append(users, c.pending...)
		}
	}

	return
}

// 获取当前这个用户连接的客服
func (c *Matcher) GetCurrentWaiter(userSocketUUID string) *string {
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

// 获取当前这个用户连接的客服
func (c *Matcher) IsUserConnectingWithWaiter(waiterUUID, userSocketUUID string) bool {
	c.RLock()
	defer c.RUnlock()

	for id, users := range c.matcher {
		if id == waiterUUID {
			for _, u := range users {
				if u == userSocketUUID {
					return true
				}
			}
		}
	}

	return false
}

// 查找一个可用的客服
func (c *Matcher) LookupWaiter() *string {
	c.RLock()
	defer c.RUnlock()
	var waiterUUID *string

	for id, users := range c.matcher {
		if len(users) < c.max {
			waiterUUID = &id
		}
	}

	return waiterUUID
}

func (c *Matcher) Connect(waiterSocketUUID string, userSocketUUID string) {
	c.RLock()
	defer c.RUnlock()

	if users, ok := c.matcher[waiterSocketUUID]; ok {
		// 如果已经连接了，那么跳过
		for _, u := range users {
			if u == userSocketUUID {
				break
			}
		}

		users = append(users, userSocketUUID)

		c.matcher[waiterSocketUUID] = users
	} else {
		c.matcher[waiterSocketUUID] = []string{userSocketUUID}
	}
}

func (c *Matcher) Disconnect(userSocketUUID string) {
	c.RLock()
	defer c.RUnlock()

	for waiterSocketUUID, users := range c.matcher {
		for index, userUUID := range users {
			if userUUID == userSocketUUID {
				// 删除
				c.matcher[waiterSocketUUID] = append(users[:index], users[index+1:]...)
			}
		}
	}
}
