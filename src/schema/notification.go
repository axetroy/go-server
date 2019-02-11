package schema

type NotificationPure struct {
	Id      string  `json:"id"`
	Author  string  `json:"author"`
	Tittle  string  `json:"tittle"`
	Content string  `json:"content"`
	Read    bool    `json:"read"` // 用户是否已读
	ReadAt  string  `json:"read"` // 用户读取的时间
	Note    *string `json:"note"`
}

type Notification struct {
	NotificationPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
