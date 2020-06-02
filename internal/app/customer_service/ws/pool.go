package ws

var (
	UserPoll   = NewPool()
	WaiterPoll = NewPool()
)

type Pool struct {
	clients   map[*Client]bool // 已连接的客户端
	Broadcast chan Message     // 广播频道
}

func (c *Pool) AddClient(client *Client) {
	c.clients[client] = true
}

func (c *Pool) GetClient(UUID string) *Client {
	for client := range c.clients {
		if client.UUID == UUID {
			return client
		}
	}
	return nil
}

func (c *Pool) RemoveClient(UUID string) {
	for client := range c.clients {
		if client.UUID == UUID {
			_ = client.Close()
			delete(c.clients, client)
		}
	}
}

func NewPool() *Pool {
	return &Pool{
		clients:   map[*Client]bool{},
		Broadcast: make(chan Message),
	}
}
