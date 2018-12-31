package model

import "time"

type NewsType string
type NewsStatus int

const (
	NewsType_News         NewsType = "news"         // 新闻咨询
	NewsType_Announcement NewsType = "announcement" // 官方公告

	NewsStatusInActive NewsStatus = -1 // 未启用的状态
	NewsStatusActive                   // 启用的状态
)

type News struct {
	Id        int64      `xorm:"pk notnull unique index" json:"id"` // 新闻公告类ID
	Author    int64      `xorm:"notnull index" json:"author"`       // 公告的作者ID
	Tittle    string     `xorm:"notnull index" json:"tittle"`       // 公告标题
	Content   string     `xorm:"notnull text" json:"content"`       // 公告内容
	Type      NewsType   `xorm:"notnull varchar(32)" json:"type"`   // 公告类型
	Tags      []string   `xorm:"notnull" json:"tags"`               // 公告的标签
	Status    NewsStatus `xorm:"notnull" json:"status"`             // 公告状态
	CreatedAt time.Time  `xorm:"created" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted" json:"deleted_at"`
}
