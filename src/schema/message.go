package schema

type MessagePure struct {
	Id      string  `json:"id"` // 消息ID
	Title   string  `json:"title"`
	Content string  `json:"content"`
	Read    bool    `json:"read"` // 用户是否已读
	Note    *string `json:"note"`
}

type Message struct {
	MessagePure
	ReadAt    *string `json:"read"` // 用户读取的时间
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}
