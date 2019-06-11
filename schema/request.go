// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

var (
	DefaultLimit = 10                // 默认只获取 10 条数据
	DefaultPage  = 0                 // 默认第 0 页
	DefaultSort  = "created_at DESC" // 默认按照创建时间排序, 只允许
	MaxLimit     = 100
)

type Query struct {
	Limit    int     `json:"limit" form:"limit"`
	Page     int     `json:"page" form:"page"`
	Sort     string  `json:"sort" form:"sort"`
	Platform *string `json:"platform" form:"platform"`
}

func (q *Query) Normalize() *Query {

	if q.Limit <= 0 {
		q.Limit = DefaultLimit // 默认查询10条
	} else if q.Limit > MaxLimit {
		q.Limit = MaxLimit
	}

	if q.Page <= 0 {
		q.Page = DefaultPage
	}

	if q.Sort == "" {
		q.Sort = DefaultSort
	}

	return q
}
