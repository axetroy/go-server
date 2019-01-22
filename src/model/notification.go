package model

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	"time"
)

type NotificationStatus int

const (
	NotificationStatusInActive NotificationStatus = -1 // 未启用的状态
	NotificationStatusActive   NotificationStatus = 0  // 启用的状态
)

type Notification struct {
	Id        string             `gorm:"primary_key;not null;unique;index;type:varchar(32)" json:"id"` // 通知ID
	Author    string             `gorm:"not null;index;type:varchar(32)" json:"Author"`                // 发布这则公告的作者
	Tittle    string             `gorm:"not null;index;type:varchar(32)" json:"tittle"`                // 公告标题
	Content   string             `gorm:"not null;text" json:"content"`                                 // 公告内容
	Status    NotificationStatus `gorm:"not null" json:"status"`                                       // 公告状态
	Note      *string            `gorm:"null;type:varchar(255)" json:"note"`                           // 这条通知的备注
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

type NotificationMark struct {
	Id        string     `gorm:"primary_key;not null;unique(uid);index;type:varchar(32)" json:"id"` // 通知ID
	Uid       string     `gorm:"notnull;index;unique(id)" json:"uid"`                               // 对应的用户ID, 联合通知ID唯一
	Read      bool       `gorm:"notnull" json:"read"`                                               // 是否已读
	ReadAt    *time.Time `gorm:"null" json:"read_at"`                                               // 阅读时间
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (news *Notification) TableName() string {
	return "notification"
}

func (news *Notification) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", util.GenerateId()); err != nil {
		return err
	}
	// 默认启用通知的状态
	if err := scope.SetColumn("status", NotificationStatusActive); err != nil {
		return err
	}
	return nil
}
