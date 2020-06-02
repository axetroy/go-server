package ws

var (
	UserPoll   = NewPool()
	WaiterPoll = NewPool()
)

type Pool struct {
	clients   map[*Client]bool // 已连接的客户端
	Broadcast chan Message     // 广播频道
}

// 添加一个连接
func (c *Pool) Add(client *Client) {
	c.clients[client] = true
}

// 获取连接
func (c *Pool) Get(UUID string) *Client {
	for client := range c.clients {
		if client.UUID == UUID {
			return client
		}
	}
	return nil
}

// 删除连接
func (c *Pool) Remove(UUID string) {
	for client := range c.clients {
		if client.UUID == UUID {
			_ = client.Close()
			delete(c.clients, client)
		}
	}
}

// 获取连接长度
func (c *Pool) Length() int {
	return len(c.clients)
}

func NewPool() *Pool {
	return &Pool{
		clients:   map[*Client]bool{},
		Broadcast: make(chan Message),
	}
}
