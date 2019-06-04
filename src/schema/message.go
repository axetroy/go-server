// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

type MessagePure struct {
	Id      string  `json:"id"`      // 消息ID
	Title   string  `json:"title"`   // 消息标题
	Content string  `json:"content"` // 消息内容
	Read    bool    `json:"read"`    // 用户是否已读
	Note    *string `json:"note"`    // 备注
}

type Message struct {
	MessagePure
	ReadAt    *string `json:"read"`       // 用户读取的时间
	CreatedAt string  `json:"created_at"` // 创建时间
	UpdatedAt string  `json:"updated_at"` // 更新时间
}
