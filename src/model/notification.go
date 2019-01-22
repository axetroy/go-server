package model

import "time"

type NotificationStatus int

const (
	NotificationStatusInActive NotificationStatus = -1 // 未启用的状态
	NotificationStatusActive   NotificationStatus = 0  // 启用的状态
)

type Notification struct {
	Id        string             `xorm:"pk notnull unique index varchar(32)" json:"id"` // 通知ID
	Tittle    string             `xorm:"notnull index varchar(32)" json:"tittle"`       // 公告标题
	Content   string             `xorm:"notnull text" json:"content"`                   // 公告内容
	Status    NotificationStatus `xorm:"notnull" json:"status"`                         // 公告状态
	Note      string             `xorm:"null varchar(255)" json:"note"`                 // 这条通知的备注
	Mark      *NotificationMark  `xorm:"'notification_mark_id'" json:"mark"`
	CreatedAt time.Time          `xorm:"created" json:"created_at"`
	UpdatedAt time.Time          `xorm:"updated" json:"updated_at"`
	DeletedAt *time.Time         `xorm:"deleted" json:"deleted_at"`
}

type NotificationMark struct {
	Id        string     `xorm:"pk notnull unique index varchar(32)" json:"id"` // 通知ID
	Nid       string     `xorm:"notnull index" json:"nid"`                      // 对应的通知ID
	Uid       string     `xorm:"notnull index unique(nid)" json:"uid"`          // 对应的用户ID, 联合通知ID唯一
	Read      bool       `xorm:"notnull" json:"read"`                           // 是否已读
	ReadAt    *time.Time `xorm:"null" json:"read_at"`                           // 阅读时间
	CreatedAt time.Time  `xorm:"created" json:"created_at"`
	UpdatedAt time.Time  `xorm:"updated" json:"updated_at"`
	DeletedAt *time.Time `xorm:"deleted" json:"deleted_at"`
}
