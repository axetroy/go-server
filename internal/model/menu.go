// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/axetroy/go-server/internal/rbac/accession"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type Menu struct {
	Id        string         `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"` // 菜单 ID
	ParentId  string         `gorm:"not null;index;type:varchar(32)" json:"parent_id"`             // 该菜单的父级 ID
	Name      string         `gorm:"not null;index;type:varchar(32)" json:"name"`                  // 菜单名
	Url       string         `gorm:"not null;index;type:varchar(255)" json:"url"`                  // 菜单链接的 URL 地址
	Icon      string         `gorm:"not null;index;type:varchar(64)" json:"icon"`                  // 菜单的图标
	Accession pq.StringArray `gorm:"not null;index;type:varchar(64)[]" json:"accession"`           // 该菜单所需要的权限
	Sort      int            `gorm:"not null;index;" json:"sort"`                                  // 菜单排序, 越大的越靠前
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

func (menu *Menu) TableName() string {
	return "menu"
}

func (menu *Menu) BeforeCreate(scope *gorm.Scope) error {
	// 检验传入的权限是否正确
	if !accession.ValidAdmin(menu.Accession) {
		return exception.InvalidParams
	}

	return scope.SetColumn("id", util.GenerateId())
}
