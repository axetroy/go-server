package notification

import "github.com/axetroy/go-server/model"

type Pure struct {
	Id      string                   `json:"id"`      // 通知ID
	Tittle  string                   `json:"tittle"`  // 公告标题
	Content string                   `json:"content"` // 公告内容
	Status  model.NotificationStatus `json:"status"`  // 公告状态
	Note    string                   `json:"note"`    // 这条通知的备注
}

type Notification struct {
	Pure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
