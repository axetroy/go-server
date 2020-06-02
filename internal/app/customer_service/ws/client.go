package ws

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
)

// 来回传输的消息体
type Message struct {
	From    string      `json:"from,omitempty"`                // 从谁发出来的
	To      string      `json:"to,omitempty"`                  // 要发送的目标 ID，只有客服才需要带 target 字段，指明发送给谁
	Type    string      `json:"type" valid:"required~请输入会话类型"` // 会话类型
	Payload interface{} `json:"payload,omitempty"`             // 本次消息的数据
}

type Client struct {
	conn *websocket.Conn // Socket 连接
	UUID string          // Socket 连接的唯一标识符
}

func NewClient(conn *websocket.Conn) *Client {
	id, err := uuid.NewRandom()

	if err != nil {
		log.Printf("%+v\n", err)
	}

	return &Client{
		conn: conn,
		UUID: id.String(),
	}
}

// 向客户端写数据
func (c *Client) WriteJSON(data Message) error {
	return c.conn.WriteJSON(data)
}

// 关闭连接
func (c *Client) Close() error {
	return c.conn.Close()
}
