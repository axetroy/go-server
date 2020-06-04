package ws

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// 来回传输的消息体
type Message struct {
	From    string      `json:"from,omitempty"`                // 从谁发出来的
	To      string      `json:"to,omitempty"`                  // 要发送的目标 ID，只有客服才需要带 target 字段，指明发送给谁
	Type    string      `json:"type" valid:"required~请输入会话类型"` // 会话类型
	Payload interface{} `json:"payload,omitempty"`             // 本次消息的数据
}

type Client struct {
	sync.RWMutex
	conn    *websocket.Conn       // Socket 连接
	UUID    string                // Socket 连接的唯一标识符
	profile *schema.ProfilePublic // 用户的身份信息，仅用于成功身份认证的连接
}

func NewClient(conn *websocket.Conn) *Client {
	id, err := uuid.NewRandom()

	if err != nil {
		log.Printf("%+v\n", err)
	}

	return &Client{
		conn:    conn,
		UUID:    id.String(),
		profile: nil,
	}
}

func (c *Client) GetProfile() *schema.ProfilePublic {
	c.Lock()
	defer c.Unlock()
	return c.profile
}

func (c *Client) UpdateProfile(profile schema.ProfilePublic) {
	c.Lock()
	defer c.Unlock()
	c.profile = &profile
}

// 向客户端写数据
func (c *Client) WriteJSON(data Message) error {
	return c.conn.WriteJSON(data)
}

func (c *Client) WriteError(err error, data Message) error {
	if e, ok := err.(exception.Error); ok {
		_ = c.WriteJSON(Message{
			Type: string(TypeResponseUserError),
			To:   c.UUID,
			Payload: map[string]interface{}{
				"message": e.Error(),
				"status":  e.Code(),
				"data":    data,
			},
		})
	} else {
		_ = c.WriteJSON(Message{
			Type: string(TypeResponseUserError),
			To:   c.UUID,
			Payload: map[string]interface{}{
				"message": err.Error(),
				"status":  exception.Unknown.Code(),
				"data":    data,
			},
		})
	}

	return c.conn.WriteJSON(data)
}

// 关闭连接
func (c *Client) Close() error {
	return c.conn.Close()
}
