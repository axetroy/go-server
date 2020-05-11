// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type ReportType string
type ReportStatus int

var (
	ReportTypeBug        ReportType = "bug"        // BUG 反馈
	ReportTypeFeature    ReportType = "feature"    // 新功能请求
	ReportTypeSuggestion ReportType = "suggestion" // 建议
	ReportTypeOther      ReportType = "other"      // 其他

	ReportStatusPending ReportStatus = 0 // 初始状态
	ReportStatusResolve ReportStatus = 1 // 已解决

	ReportTypes = []ReportType{ReportTypeBug, ReportTypeFeature, ReportTypeSuggestion, ReportTypeOther}
)

// 检验是否是有效的报错类型
func IsValidReportType(t ReportType) bool {
	for _, v := range ReportTypes {
		if v == t {
			return true
		}
	}
	return false
}

type Report struct {
	Id          string         `gorm:"primary_key;unique;not null;index;type:varchar(32)" json:"id"` // 反馈ID
	Uid         string         `gorm:"not null;index;type:varchar(32)" json:"uid"`                   // 反馈的作者ID
	Title       string         `gorm:"not null;index;type:varchar(32)" json:"title"`                 // 反馈标题
	Content     string         `gorm:"not null;type:text" json:"content"`                            // 反馈内容
	Type        ReportType     `gorm:"not null;type:varchar(32)" json:"type"`                        // 反馈类型
	Status      ReportStatus   `gorm:"not null;" json:"status"`                                      // 当前报告的处理状态
	Screenshots pq.StringArray `gorm:"type:varchar(256)[]" json:"screenshots"`                       // 反馈的截图
	Locked      bool           `gorm:"not null;" json:"locked"`                                      // 是否已锁定，锁定之后用户不能再更改状态
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

func (report *Report) TableName() string {
	return "report"
}

func (report *Report) BeforeCreate(scope *gorm.Scope) (err error) {
	err = scope.SetColumn("id", util.GenerateId())
	err = scope.SetColumn("status", ReportStatusPending)
	return
}
