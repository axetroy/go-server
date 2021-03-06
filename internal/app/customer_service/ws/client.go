// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package ws

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/schema"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

// 来回传输的消息体
type Message struct {
	Id      string      `json:"id,omitempty" validate:"omitempty" comment:"消息 ID"`        // 每条消息的 ID，在写入数据库之后会有
	From    string      `json:"from,omitempty" validate:"omitempty,uuid" comment:"发送者"`   // 从谁发出来的
	To      string      `json:"to,omitempty" validate:"omitempty,uuid" comment:"发送目标"`    // 要发送的目标 ID，只有客服才需要带 target 字段，指明发送给谁
	Type    string      `json:"type" validate:"required" comment:"会话类型"`                  // 会话类型
	Payload interface{} `json:"payload,omitempty" validate:"omitempty" comment:"消息数据"`    // 本次消息的数据
	Date    string      `json:"date,omitempty" validate:"omitempty" comment:"时间戳"`        // 消息的时间
	OpID    *string     `json:"op_id,omitempty" validate:"omitempty,uuid" comment:"操作ID"` // 客户端发送的 UUID 用于标记该消息，服务端应该返回相同的 op_id，主要用于消息回执
}

type Client struct {
	sync.RWMutex
	conn            *websocket.Conn       // Socket 连接
	UUID            string                // Socket 连接的唯一标识符
	profile         *schema.ProfilePublic // 用户的身份信息，仅用于成功身份认证的连接
	LatestReceiveAt time.Time             // 最近接收到的消息的时间，用于判断用户是否空闲
	Closed          bool                  // 连接是否已关闭
	Ready           bool                  // 该客户端是否已准备就绪，给客服端用的，ready  = true 的时候系统才会分配用户
}

func NewClient(conn *websocket.Conn) *Client {
	id, err := uuid.NewRandom()

	if err != nil {
		log.Printf("%+v\n", err)
	}

	return &Client{
		conn:            conn,
		UUID:            id.String(),
		profile:         nil,
		LatestReceiveAt: time.Now(),
	}
}

func (c *Client) GetProfile() *schema.ProfilePublic {
	c.Lock()
	defer c.Unlock()
	return c.profile
}

func (c *Client) RegenerateUUID() {
	c.Lock()
	defer c.Unlock()
	id, _ := uuid.NewRandom()

	c.UUID = id.String()
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
		if e1 := c.WriteJSON(Message{
			Type: string(TypeResponseUserError),
			To:   c.UUID,
			Payload: schema.Response{
				Message: e.Error(),
				Status:  e.Code(),
				Data:    data,
			},
			Date: time.Now().Format(time.RFC3339Nano),
			OpID: data.OpID,
		}); e1 != nil {
			return e1
		}
	} else {
		if e2 := c.WriteJSON(Message{
			Type: string(TypeResponseUserError),
			To:   c.UUID,
			Payload: schema.Response{
				Message: e.Error(),
				Status:  exception.Unknown.Code(),
				Data:    data,
			},
			Date: time.Now().Format(time.RFC3339Nano),
			OpID: data.OpID,
		}); e2 != nil {
			return e2
		}
	}

	return nil
}

// 关闭连接
func (c *Client) Close() error {
	c.Lock()
	defer c.Unlock()
	c.Closed = true
	return c.conn.Close()
}

func (c *Client) SetReady(ready bool) {
	c.Lock()
	defer c.Unlock()
	c.Ready = ready
}
