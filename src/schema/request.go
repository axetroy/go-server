// Copyright 2019 Axetroy. All rights reserved. MIT license.
package schema

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"regexp"
	"strings"
)

type Order string

type Query struct {
	Limit    int     `json:"limit" form:"limit"`
	Page     int     `json:"page" form:"page"`
	Sort     string  `json:"sort" form:"sort"`
	Platform *string `json:"platform" form:"platform"`
}

type Sort struct {
	Field string `json:"field"` // 排序的字段
	Order Order  `json:"order"` // 字段的排序方向
}

var (
	DefaultLimit       = 10            // 默认只获取 10 条数据
	DefaultPage        = 0             // 默认第 0 页
	DefaultSort        = "-created_at" // 默认按照创建时间排序
	MaxLimit           = 100           // 最大的查询数量，100 条 防止查询数据过大拖慢服务端性能
	OrderAsc     Order = "ASC"         // 排序方式，正序
	OrderDesc    Order = "DESC"        // 排序方式，倒序
	ascReg             = regexp.MustCompile("^\\+")
	descReg            = regexp.MustCompile("^-")
)

func NewQuery() *Query {
	q := Query{}

	q.Normalize()

	return &q
}

func (q *Query) FormatSort() (fields []Sort) {
	arr := strings.Split(q.Sort, ",")

	for _, field := range arr {
		s := strings.Split(field, "")

		switch s[0] {

		case "-":
			fields = append(fields, Sort{
				Field: descReg.ReplaceAllString(field, ""),
				Order: OrderDesc,
			})
		default:
			fields = append(fields, Sort{
				Field: ascReg.ReplaceAllString(field, ""),
				Order: OrderAsc,
			})
		}
	}

	return
}

func (q *Query) Order(db *gorm.DB) *gorm.DB {
	sorts := q.FormatSort()

	for _, field := range sorts {
		db = db.Order(fmt.Sprintf("%s %s", field.Field, field.Order))
	}

	return db
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
