package schema

import "github.com/axetroy/go-server/src/model"

type BannerPure struct {
	Id          string               `json:"id"`           // 地址ID
	Image       string               `json:"image"`        // 图片 URL
	Href        string               `json:"href"`         // 点击图片跳转 URL
	Platform    model.BannerPlatform `json:"platform"`     // 平台
	Description *string              `json:"description"`  // 描述
	Priority    *string              `json:"priority"`     // 优先级，用于排序
	Identifier  *string              `json:"identifier"`   // APP 跳转标识符
	FallbackUrl *string              `json:"fallback_url"` // APP 跳转标识符的备选方案
}

type Banner struct {
	BannerPure
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
