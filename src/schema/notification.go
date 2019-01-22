package schema

type NotificationPure struct {
	Id      string  `json:"id"`
	Author  string  `json:"author"`
	Tittle  string  `json:"tittle"`
	Content string  `json:"content"`
	Note    *string `json:"note"`
}

type Notification struct {
	NotificationPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
