// Copyright 2019 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type HelpStatus int
type HelpType string

const (
	HelpStatusInActive HelpStatus = -1        // 未启用的状态
	HelpStatusActive   HelpStatus = 1         // 启用的状态
	HelpTypeArticle    HelpType   = "article" // 普通文章
	HelpTypeClass      HelpType   = "class"   // 分类
)

type Help struct {
	Id        string         `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"` // 帮助文章ID
	Title     string         `gorm:"not null;index;type:varchar(32)" json:"title"`                 // 帮助文章标题
	Content   string         `gorm:"not null;type:text" json:"content"`                            // 帮助文章内容
	Tags      pq.StringArray `gorm:"type:varchar(32)[]" json:"tags"`                               // 帮助文章的标签
	Status    HelpStatus     `gorm:"not null;type:integer" json:"status"`                          // 帮助文章状态
	Type      HelpType       `gorm:"not null;type:varchar(32)" json:"type"`                        // 帮助文章的类型
	ParentId  *string        `gorm:"null;index;type:varchar(32)" json:"parent_id"`                 // 父级 ID，如果有的话                                    // 父级 ID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (news *Help) TableName() string {
	return "help"
}

func (news *Help) BeforeCreate(scope *gorm.Scope) error {
	if err := scope.SetColumn("id", util.GenerateId()); err != nil {
		return err
	}

	return nil
}
