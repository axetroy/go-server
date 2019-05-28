package model

import (
	"github.com/axetroy/go-server/src/util"
	"github.com/jinzhu/gorm"
	"github.com/lib/pq"
	"time"
)

type ReportType string

var (
	Bug        ReportType = "bug"        // BUG 反馈
	Feature    ReportType = "feature"    // 新功能请求
	Suggestion ReportType = "suggestion" // 建议
	Other      ReportType = "other"      // 其他

	ReportTypes = []ReportType{Bug, Feature, Suggestion, Other}
)

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
	Screenshots pq.StringArray `gorm:"type:varchar(256)[]" json:"screenshots"`                       // 反馈的截图
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
}

func (report *Report) TableName() string {
	return "report"
}

func (report *Report) BeforeCreate(scope *gorm.Scope) error {
	return scope.SetColumn("id", util.GenerateId())
}
