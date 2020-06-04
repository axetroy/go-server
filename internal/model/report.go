// Copyright 2019-2020 Axetroy. All rights reserved. MIT license.
package model

import (
	"github.com/axetroy/go-server/internal/library/exception"
	"github.com/axetroy/go-server/internal/library/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type ReportType string
type ReportStatus int
type ReportTypeDetail struct {
	Type        ReportType `json:"type"`
	Description string     `json:"description"`
}

var (
	ReportTypeBug        ReportType = "bug"        // BUG 反馈
	ReportTypeFeature    ReportType = "feature"    // 新功能请求
	ReportTypeSuggestion ReportType = "suggestion" // 建议
	ReportTypeOther      ReportType = "other"      // 其他
	ReportTypes                     = []ReportTypeDetail{
		{
			Type:        ReportTypeBug,
			Description: "BUG 反馈",
		},
		{
			Type:        ReportTypeFeature,
			Description: "新功能请求",
		},
		{
			Type:        ReportTypeSuggestion,
			Description: "建议",
		},
		{
			Type:        ReportTypeOther,
			Description: "其他",
		},
	}

	ReportStatusPending ReportStatus = 0 // 初始状态
	ReportStatusResolve ReportStatus = 1 // 已解决
	ReportStatuses                   = []ReportStatus{ReportStatusPending, ReportStatusResolve}
)

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

func (r *Report) TableName() string {
	return "report"
}

func (r *Report) IsValidType() bool {
	for _, v := range ReportTypes {
		if v.Type == r.Type {
			return true
		}
	}
	return false
}

func (r *Report) IsValidStatus() bool {
	for _, v := range ReportStatuses {
		if v == r.Status {
			return true
		}
	}
	return false
}

func (r *Report) BeforeCreate(scope *gorm.Scope) (err error) {
	// 校验 type 是否正确
	if r.IsValidType() == false {
		return exception.InvalidParams.New("无效的类型")
	}
	if r.IsValidStatus() == false {
		return exception.InvalidParams.New("无效的状态")
	}

	if err := scope.SetColumn("id", util.GenerateId()); err != nil {
		return err
	}
	if err := scope.SetColumn("status", ReportStatusPending); err != nil {
		return err
	}
	return
}
